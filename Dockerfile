FROM alpine:3.4
MAINTAINER Johnny Bergström <johnny@joonix.se>

RUN apk --update add ca-certificates

RUN mkdir /app
ADD joxy /app

WORKDIR /app
EXPOSE 443
ENTRYPOINT ["/app/joxy"]