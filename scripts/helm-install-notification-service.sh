#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# Get the project root directory (parent of scripts/)
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

helm install notification-service "$PROJECT_ROOT/helm/notification-service" \
  -f "$PROJECT_ROOT/helm/notification-service/values.yaml" \
  -f "$PROJECT_ROOT/helm/notification-service/values-secret.yaml" \
  --debug