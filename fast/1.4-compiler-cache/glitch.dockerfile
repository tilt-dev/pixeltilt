FROM golang:alpine
WORKDIR /app
COPY vendor vendor
COPY render/api render/api
COPY go.mod ./
COPY glitch glitch
RUN GOCACHE=/app/glitch/cache CGO_ENABLED=0 go build -mod=vendor -ldflags '-w' glitch/main.go
CMD ["./main"]