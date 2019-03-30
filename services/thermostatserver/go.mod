module github.com/mrbenshef/SmartHomeAdapters/thermostatserver

require (
	github.com/golang/protobuf v1.3.1
	github.com/lib/pq v1.0.0
	github.com/mrbenshef/SmartHomeAdapters/microservice v0.0.0-20190213130345-d379925727ab
	google.golang.org/grpc v1.19.0
)

replace github.com/mrbenshef/SmartHomeAdapters/microservice => ../microservice
