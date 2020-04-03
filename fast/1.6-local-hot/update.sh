#!/bin/bash

sed -re 's/time.Unix\(0, ([0-9]*)\)/time.Unix\(0, '"$(($(date +%s%N)))"'\)/' -i glitch/main.go

CGO_ENABLED=0 go build -o glitch/main -ldflags '-w' glitch/main.go