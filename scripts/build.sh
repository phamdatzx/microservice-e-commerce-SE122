#!/bin/bash

set -e

read -p "Enter the version: " VERSION
read -p "Enter the service names (comma separated): " SERVICE_NAMES
read -p "Enter the container registry: " CONTAINER_REGISTRY

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