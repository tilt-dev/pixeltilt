#!/bin/bash

sed -re 's/time.Unix\(0, ([0-9]*)\)/time.Unix\(0, '"$(($(date +%s%N)))"'\)/' -i glitch/main.go

CGO_ENABLED=0 go build -o bigbinary -ldflags '-w' glitch/main.go

upx bigbinary -1

mv bigbinary glitch/main