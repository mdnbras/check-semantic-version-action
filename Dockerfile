FROM golang:1.21-alpine3.18 as builder

RUN mkdir /go/src/checkSemanticVersion

WORKDIR /go/src/checkSemanticVersion

RUN apk add nano zip git

COPY ./ ./

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o executable

RUN zip executable.zip executable