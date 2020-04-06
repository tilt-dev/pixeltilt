#!/bin/bash

read -r -d '' VAR <<- EOM
apiVersion: apps/v1
kind: Deployment
metadata:
  name: $name
  labels:
    app: $name
spec:
  selector:
    matchLabels:
      app: $name
  template:
    metadata:
      labels:
        app: $name
        tier: web
    spec:
      containers:
      - name: $name
        image: $name
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: $name
  labels:
    app: $name
spec:
  ports:
    - port: 8080
      targetPort: 8080
      protocol: TCP
  selector:
    app: $name
EOM
echo "$VAR"




