#!/bin/sh
echo Compiling...
go build -mod=vendor -ldflags '-w' /app/muxer/main.go
./main