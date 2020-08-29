FROM alpine:3.12
ADD core/bin/digger /usr/bin/
ADD ui/dist /var/www/html
RUN apk add tzdata libc6-compat
WORKDIR /var/www/html
