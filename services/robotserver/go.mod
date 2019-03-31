module github.com/mrbenshef/SmartHomeAdapters/robotserver

require (
	github.com/golang/protobuf v1.3.1
	github.com/gorilla/websocket v1.4.0
	github.com/julienschmidt/httprouter v1.2.0
	github.com/mrbenshef/SmartHomeAdapters/microservice v0.0.0
	google.golang.org/grpc v1.19.0
)

replace github.com/mrbenshef/SmartHomeAdapters/microservice => ../microservice
