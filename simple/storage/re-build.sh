#!/bin/sh
echo Compiling...
go build -mod=vendor -ldflags '-w' /app/storage/main.go
./main