#!/bin/sh
echo Compiling...
go build -mod=vendor -ldflags '-w' /app/rectangler/main.go
./main