#!/bin/bash

read -r -d '' VAR <<- EOM
FROM alpine
RUN apk add entr
RUN mkdir /app
WORKDIR /app
ADD $files ./
CMD ls main | entr -n -r ./main
EOM
echo "$VAR"
