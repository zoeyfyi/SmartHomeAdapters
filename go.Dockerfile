# build stage
FROM golang AS build-env

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
COPY thermostatserver/go.mod thermostatserver/go.mod
COPY thermostatserver/go.sum thermostatserver/go.sum
COPY account-app/go.mod account-app/go.mod
COPY account-app/go.sum account-app/go.sum
COPY account-app/static/login.html account-app/static/login.html
COPY account-app/static/consent.html account-app/static/consent.html


RUN cd clientserver && go mod download
RUN cd infoserver && go mod download
RUN cd robotserver && go mod download
RUN cd switchserver && go mod download
RUN cd userserver && go mod download
RUN cd thermostatserver && go mod download
RUN cd account-app && go mod download

# build server
ARG SERVICE

COPY . .
RUN cd $SERVICE && go build -o server
RUN cd $SERVICE && mkdir /app && mv server /app/server

RUN if [ "$SERVICE" = "account-app" ] ; then cp -r account-app/static /app ; fi
RUN ls /app/static/login.html
RUN cat /app/static/login.html
# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /app/server /app/
ENTRYPOINT ./server
