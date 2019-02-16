module github.com/mrbenshef/SmartHomeAdapters/thermostatserver

require (
	github.com/golang/protobuf v1.2.0
	github.com/lib/pq v1.0.0
	github.com/mrbenshef/SmartHomeAdapters/robotserver v0.0.0-20190213130345-d379925727ab
	golang.org/x/net v0.0.0-20180826012351-8a410e7b638d
	google.golang.org/grpc v1.18.0
)

replace github.com/mrbenshef/SmartHomeAdapters/robotserver => ../robotserver
