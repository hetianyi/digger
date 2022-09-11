FROM golang:1.12-alpine
ADD . /digger
RUN sed -i "s@http://dl-cdn.alpinelinux.org/@https://repo.huaweicloud.com/@g" /etc/apk/repositories && \
    apk add git gcc musl-dev pkgconfig tzdata libc6-compat && \
    cd /digger && \
    sh build-core.sh

FROM alpine:3.12
COPY --from=0 /digger/core/bin/digger /usr/bin/
RUN sed -i "s@http://dl-cdn.alpinelinux.org/@https://repo.huaweicloud.com/@g" /etc/apk/repositories && \
    apk add tzdata libc6-compat
ADD ui/dist /var/www/html
WORKDIR /var/www/html