#!/bin/sh
echo Compiling...
go build -mod=vendor -ldflags '-w' /app/red/main.go
./main