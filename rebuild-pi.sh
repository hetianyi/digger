sh build-all-arm32v7.sh
docker build -t hehety/digger:arm32v7-latest . -f Dockerfile-arm32v7
docker push hehety/digger:arm32v7-latest
