#!/bin/sh

echo "build ui..."
docker run --rm \
  -v $PWD:/app \
  -w /app \
  node:9.7.0 \
  sh -c 'cd ui && npm i && npm run build'

build_status=$?
if [ $build_status != 0 ]; then
  echo "build failed"
  exit build_status
fi
echo "build ui success"
