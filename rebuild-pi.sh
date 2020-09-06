sh build-all-arm32v7.sh
tag=RELEASE.$(date "+%Y-%m-%d_%H-%M-%S")
docker build -t hehety/digger:${tag}-arm . -f Dockerfile-arm32v7
docker tag hehety/digger:${tag}-arm hehety/digger:latest-arm
docker push hehety/digger:${tag}-arm
docker push hehety/digger:latest-arm
docker rmi hehety/digger:${tag}-arm
