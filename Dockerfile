# syntax = docker/dockerfile:1

FROM golang:1.19-alpine AS build
WORKDIR /src
RUN apk add --no-cache file git
ENV GOMODCACHE /root/.cache/gocache
RUN --mount=target=. --mount=target=/root/.cache,type=cache \
    CGO_ENABLED=0 go build -o /bin/cloak -ldflags '-s -d -w' ./cmd/cloak

FROM alpine:3.16
RUN apk add -U --no-cache ca-certificates
COPY --from=docker:20.10.17-cli-alpine3.16 /usr/local/bin/docker /usr/local/bin/docker
COPY --from=build /bin/cloak /bin/cloak
ENTRYPOINT ["/bin/cloak"]
