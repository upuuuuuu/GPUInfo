# FROM golang:1.18 as builder

# WORKDIR /app

# COPY . .
# RUN go build -o ingress-manager main.go

FROM ubuntu:20.04

WORKDIR /app

COPY  gpu /app/

CMD [ "./gpu" ]
