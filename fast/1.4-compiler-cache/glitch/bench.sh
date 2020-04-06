#!/bin/bash
rm -rf cache2
mkdir cache2

printf "Without compiler cache:\t"
GOCACHE=$(pwd)/cache2 CGO_ENABLED=0 /usr/bin/time -f "%Es" go build -o binary

printf "With compiler cache:\t"
GOCACHE=$(pwd)/cache2 CGO_ENABLED=0 /usr/bin/time -f "%Es" go build -o binary

rm binary
