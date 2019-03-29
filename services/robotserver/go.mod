module github.com/mrbenshef/SmartHomeAdapters/robotserver

require (
	github.com/golang/protobuf v1.3.1
	github.com/google/pprof v0.0.0-20190109223431-e84dfd68c163 // indirect
	github.com/gorilla/websocket v1.4.0
	github.com/ianlancetaylor/demangle v0.0.0-20181102032728-5e5cf60278f6 // indirect
	github.com/julienschmidt/httprouter v1.2.0
	github.com/mrbenshef/SmartHomeAdapters/microservice v0.0.0
	golang.org/x/arch v0.0.0-20181203225421-5a4828bb7045 // indirect
	golang.org/x/net v0.0.0-20190318221613-d196dffd7c2b
	google.golang.org/grpc v1.19.0
)

replace github.com/mrbenshef/SmartHomeAdapters/microservice => ../microservice
