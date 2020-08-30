#!/bin/sh
cd core
echo "build core..."
docker run --rm \
  -v ~/go:/go \
  -v $PWD:/app \
  -w /app \
  -e GOPROXY=https://goproxy.io  \
  golang:1.12-stretch \
  go build -o bin/digger

build_status=$?
if [ $build_status != 0 ]; then
  echo "build failed"
  exit build_status
fi
echo "build core success"
