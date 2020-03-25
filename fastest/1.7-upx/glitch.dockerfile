FROM alpine
RUN apk add entr
RUN mkdir /app 
WORKDIR /app 
ADD glitch/main ./
CMD ls main | entr -r ./main