FROM golang:1.19 AS builder

WORKDIR /build/cmd

COPY ./pow-server/ ../

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./cmd/main ./
