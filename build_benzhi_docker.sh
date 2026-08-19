#!/bin/bash
set -e

IMAGE_NAME="tailings-dam"
IMAGE_TAG="latest"

echo "=============================="
echo "Building ${IMAGE_NAME}:${IMAGE_TAG}"
echo "=============================="

docker build -f benzhi.Dockerfile -t ${IMAGE_NAME}:${IMAGE_TAG} .

echo ""
echo "=============================="
echo "Build complete!"
echo "=============================="
echo ""
echo "Run with:"
echo "  docker run -p 8080:8080 -v tailings-dam-data:/data ${IMAGE_NAME}:${IMAGE_TAG}"
echo ""
echo "API available at: http://localhost:8080/api/dashboard"
