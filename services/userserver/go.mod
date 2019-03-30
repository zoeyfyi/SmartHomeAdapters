module github.com/mrbenshef/SmartHomeAdapters/userserver

require (
	github.com/golang/protobuf v1.3.1
	github.com/lib/pq v1.0.0
	github.com/mrbenshef/SmartHomeAdapters/microservice v0.0.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2
	google.golang.org/grpc v1.19.0
)

replace github.com/mrbenshef/SmartHomeAdapters/microservice => ../microservice
