#!/bin/bash

set -euo pipefail

if [ $# -lt 1 ] || [ $# -gt 2 ]; then
  echo "Usage: $0 <user@host> [remote_app_dir]"
  echo "Example: $0 ubuntu@your-vps /app"
  exit 1
fi

REMOTE="$1"
REMOTE_APP_DIR="${2:-/app}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
COMPOSE_FILE="$PROJECT_ROOT/docker-compose.yaml"

if [ ! -f "$COMPOSE_FILE" ]; then
  echo "Compose file not found: $COMPOSE_FILE"
  exit 1
fi

cd "$PROJECT_ROOT"

mapfile -t ENV_FILES < <(
  awk '
    /^[[:space:]]*-[[:space:]]*/ {
      value = $0
      sub(/^[[:space:]]*-[[:space:]]*/, "", value)
      gsub(/["'\''[:space:]]/, "", value)
      if (value ~ /^(\.\/)?services\/.*\/\.env$/) {
        sub(/^\.\//, "", value)
        print value
      }
    }
  ' "$COMPOSE_FILE" | sort -u
)

if [ "${#ENV_FILES[@]}" -eq 0 ]; then
  echo "No service .env files found in $COMPOSE_FILE"
  exit 1
fi

for env_file in "${ENV_FILES[@]}"; do
  if [ ! -f "$env_file" ]; then
    echo "Missing env file: $env_file"
    exit 1
  fi
done

echo "Syncing env files to $REMOTE:$REMOTE_APP_DIR"
printf ' - %s\n' "${ENV_FILES[@]}"

REMOTE_APP_DIR_QUOTED="$(printf '%q' "$REMOTE_APP_DIR")"
REMOTE_ENV_DIRS=()

for env_file in "${ENV_FILES[@]}"; do
  REMOTE_ENV_DIRS+=("$(printf '%q' "$REMOTE_APP_DIR/$(dirname "$env_file")")")
done

ssh "$REMOTE" "mkdir -p $REMOTE_APP_DIR_QUOTED ${REMOTE_ENV_DIRS[*]}"

tar -czf - "${ENV_FILES[@]}" | ssh "$REMOTE" \
  "tar -xzf - -C $REMOTE_APP_DIR_QUOTED"

echo "Env files synced successfully."
