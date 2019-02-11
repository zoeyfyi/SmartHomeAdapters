# build stage
FROM golang AS build-env

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# cache dependencies
COPY go.mod /src/go.mod
COPY go.sum /src/go.sum
RUN cd /src && go mod download

# build server
COPY . /src
RUN cd /src && go generate 
RUN cd /src && go build -o server

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/server /app/
ENTRYPOINT ./server