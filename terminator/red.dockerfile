FROM golang:1.13.6-alpine

RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY render render
COPY red red

RUN GO111MODULE=on go build red/main.go

EXPOSE 8085

CMD ["./main"]
