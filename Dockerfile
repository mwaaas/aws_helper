FROM golang:1.11.5-alpine3.8 as build-env

RUN apk add --update \
    && apk add ca-certificates git gcc musl-dev \
    && mkdir -p ./build

WORKDIR /usr/src/app
ENV Port 80
COPY . .

RUN  go build  -o dist/exporter

CMD go run *.go


# running container
FROM alpine:3.8
RUN apk add --update \
      ca-certificates
ENV Port 80
WORKDIR /app
COPY --from=build-env /usr/src/app/dist/exporter /app/
ENTRYPOINT ./exporter