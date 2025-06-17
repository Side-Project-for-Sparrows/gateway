#!/bin/bash

set -e

IMAGE_NAME="gateway-service:stg"

echo "🛠 Docker 이미지 빌드 중..."
docker build --no-cache -t $IMAGE_NAME .

echo "🚀 컨테이너 실행 중..."
docker run --rm -p 7080:7080 $IMAGE_NAME

