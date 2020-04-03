#!/bin/sh
echo Compiling...
go build -mod=vendor -ldflags '-w' glitch/main.go
./main