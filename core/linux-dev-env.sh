#!/bin/bash

docker rm -f postgres
docker run -d \
  -p 5432:5432 \
  --name postgres \
  --restart always \
  -e TZ=Asia/Shanghai \
  -e POSTGRES_USER=digger \
  -e POSTGRES_PASSWORD=123456 \
  -e POSTGRES_DB=digger \
  -e PGDATA=/var/lib/postgresql/data \
  -v postgres:/var/lib/postgresql/data \
  -v $PWD/../db:/docker-entrypoint-initdb.d:ro \
  postgres:9.6 \
  --timezone=PRC


docker run -f redis
docker run -d \
  -p 6379:6379 \
  --name redis \
  -e TZ=Asia/Shanghai \
  -v redis:/data \
  redis:5.0.9 \
  redis-server --bind 0.0.0.0 --protected-mode no --appendonly yes --requirepass 123456










