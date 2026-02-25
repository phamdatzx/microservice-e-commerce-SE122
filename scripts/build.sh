#!/bin/bash

set -e

VERSION="$1"
SERVICE_NAMES="$2"
CONTAINER_REGISTRY="$3"

if [ $# -lt 3 ]; then
  echo "Usage: $0 <version> <service1,service2> <container_registry>"
  exit 1
fi

# Split service names into an array
IFS=',' read -ra SERVICES <<< "$SERVICE_NAMES"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Build and push each service
for SERVICE in "${SERVICES[@]}"; do
  echo "Building and pushing $SERVICE..."
  docker build -t "$CONTAINER_REGISTRY/$SERVICE:$VERSION" "$SCRIPT_DIR/../services/$SERVICE"
  docker push "$CONTAINER_REGISTRY/$SERVICE:$VERSION"
done

echo "All services have been built and pushed successfully!"