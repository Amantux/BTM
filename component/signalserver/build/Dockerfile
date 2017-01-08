FROM golang:1.7.1

COPY SignalServer.go /usr/src/

WORKDIR /usr/src

ENV GOOS=linux
ENV GOARCH=arm

RUN go build -o /go/bin/SignalServer


