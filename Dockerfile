# syntax=docker/dockerfile:1

FROM golang:latest as base

COPY . /app

WORKDIR /app

RUN go mod download

RUN go build -o /bin/main -v ./cmd/main

EXPOSE 80

CMD [ "/bin/main" ]