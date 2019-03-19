# build stage
FROM smarthomeadapters/microservice AS build-env

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o server

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/server /app/
ENTRYPOINT ./server
