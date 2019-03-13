# build stage
FROM golang AS build-env

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64
WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o server

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/server /app/
COPY --from=build-env /app/static /app/static
ENTRYPOINT ./server
