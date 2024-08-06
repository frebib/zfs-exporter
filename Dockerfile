FROM golang:alpine3.20 AS build

RUN apk add \
        gcc \
        git \
        libc-dev \
        zfs-dev

WORKDIR /build

ADD . .
RUN go build -v -trimpath


FROM alpine:3.20

RUN apk add zfs-libs
COPY --from=build /build/zfs-exporter /usr/bin/

EXPOSE 9254
ENTRYPOINT ["/usr/bin/zfs-exporter"]
