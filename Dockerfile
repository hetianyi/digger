FROM golang:1.12-alpine
RUN apk add git && \
    git clone https://github.com/hetianyi/digger.git && \
    cd /digger && \
    git checkout buildx && \
    sh build-core.sh

FROM alpine:3.12
COPY --from=0 /digger/core/bin/digger /usr/bin/
RUN apk add tzdata && apk add libc6-compat
ADD ui/dist /var/www/html
WORKDIR /var/www/html