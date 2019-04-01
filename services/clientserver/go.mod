module github.com/mrbenshef/SmartHomeAdapters/clientserver

require (
	github.com/golang/protobuf v1.3.1
	github.com/julienschmidt/httprouter v1.2.0
	github.com/mrbenshef/SmartHomeAdapters/microservice v0.0.0
	google.golang.org/genproto v0.0.0-20180831171423-11092d34479b // indirect
	google.golang.org/grpc v1.19.0
)

replace github.com/mrbenshef/SmartHomeAdapters/microservice => ../microservice
