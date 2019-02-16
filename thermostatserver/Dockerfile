FROM golang

ENV GO111MODULE=on

WORKDIR /app
COPY . .

RUN if [ -d "vendor" ]; then CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -o server; else CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server; fi

CMD ["/app/server"]