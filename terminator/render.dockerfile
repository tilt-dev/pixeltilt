FROM golang:1.13.6-alpine

RUN apk update && apk upgrade && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY frontend frontend
COPY storage storage
COPY render render

RUN GO111MODULE=on go build render/main.go

EXPOSE 8084

CMD ["./main"]
