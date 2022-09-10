FROM golang:1.12-alpine
ADD . /digger
RUN cd digger && \
    sh build-core.sh

FROM alpine:3.12
COPY --from=0 /digger/core/bin/digger /usr/bin/
RUN apk add tzdata && apk add libc6-compat
ADD ui/dist /var/www/html
WORKDIR /var/www/html