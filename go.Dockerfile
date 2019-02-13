# build stage
FROM golang AS build-env

ARG SERVICE

ENV GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# cache dependencies
COPY clientserver/go.mod clientserver/go.mod
COPY clientserver/go.sum clientserver/go.sum
COPY infoserver/go.mod infoserver/go.mod
COPY infoserver/go.sum infoserver/go.sum
COPY robotserver/go.mod robotserver/go.mod
COPY robotserver/go.sum robotserver/go.sum
COPY switchserver/go.mod switchserver/go.mod
COPY switchserver/go.sum switchserver/go.sum
COPY userserver/go.mod userserver/go.mod
COPY userserver/go.sum userserver/go.sum

RUN cd $SERVICE && go mod download

# build server
COPY . .
RUN cd $SERVICE && go build -o server
RUN cd $SERVICE && mkdir /app && mv server /app/server

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/server /app/
ENTRYPOINT ./server