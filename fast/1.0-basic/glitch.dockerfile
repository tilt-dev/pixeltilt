FROM golang:alpine
WORKDIR /app
COPY render/api render/api
COPY go.mod ./
COPY glitch glitch
RUN go build glitch/main.go
CMD ["./main"]