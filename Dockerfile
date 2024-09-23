FROM golang:1.20-alpine3.16 as build
RUN apk add build-base git openssh-client openssl-dev librdkafka-dev librdkafka pkgconf

RUN mkdir /app

COPY . /app/

WORKDIR /app

RUN go build -tags musl -o /ingestion-films cmd/main.go

FROM alpine:3.19
LABEL MAINTAINER="xochilpili <xochilpili@gmail.com>"
RUN apk add --no-cache ca-certificates
COPY --from=build /app /app
ENTRYPOINT /app