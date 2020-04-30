#!/bin/sh
echo Compiling...
go build -mod=vendor -ldflags '-w' /app/glitch/main.go
./main