FROM golang:alpine
WORKDIR /app
RUN apk add entr
COPY vendor vendor
COPY render/api render/api
COPY go.mod compile.sh ./
COPY glitch glitch
CMD ls glitch/*.go | entr -r ./compile.sh