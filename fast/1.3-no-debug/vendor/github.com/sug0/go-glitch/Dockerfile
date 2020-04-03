#
# Build
#
FROM golang:1.9 as builder

WORKDIR /go/src/github.com/sugoiuguu/go-glitch
COPY ./ ./

WORKDIR example

# deps
RUN go get -d -v

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /glitch glitch.go

#
# Final image
#
FROM alpine

COPY --from=builder /glitch /usr/bin/

WORKDIR /root
ENTRYPOINT ["glitch"]
