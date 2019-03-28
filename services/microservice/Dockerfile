# build stage
FROM golang

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /microservice

COPY go.mod .
RUN go mod download

COPY . .
