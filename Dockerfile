FROM golang:alpine

# Set destination for COPY
WORKDIR /usr/src/app

# Allow for running `go test`
RUN apk add build-base

# Download Go modules
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY *.go ./
