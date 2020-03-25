FROM golang:1.13.6-alpine

RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY muxer muxer
COPY storage storage
COPY render/api render/api

RUN GO111MODULE=on go build muxer/main.go

EXPOSE 8080

CMD ["./main"]
