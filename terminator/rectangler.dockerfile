FROM golang:1.13.6-alpine

RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY render/api render/api
COPY render/Inconsolata-Bold.ttf render/Inconsolata-Bold.ttf
COPY rectangler rectangler

RUN GO111MODULE=on go build rectangler/main.go

EXPOSE 8085

CMD ["./main"]
