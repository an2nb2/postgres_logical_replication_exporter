FROM golang:1.19-alpine AS build

ADD . /src

WORKDIR /src

ENV CGO_ENABLED=0

RUN go build -o ./bin/exporter -mod=readonly


FROM alpine:3.16

COPY --from=build /src/bin/exporter /

RUN apk add --no-cache libc6-compat

ENTRYPOINT ["/exporter"]
