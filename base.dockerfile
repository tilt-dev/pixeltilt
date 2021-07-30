FROM golang:1.16-alpine

WORKDIR /app

ENV GOMODCACHE=/cache/gomod
ENV GOCACHE=/cache/gobuild

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/cache/gomod \
    go mod download
