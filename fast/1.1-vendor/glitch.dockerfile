FROM golang:alpine
WORKDIR /app
COPY vendor vendor
COPY render/api render/api
COPY go.mod ./
COPY glitch glitch
RUN go build -mod=vendor glitch/main.go
CMD ["./main"]