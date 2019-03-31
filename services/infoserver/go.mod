module github.com/mrbenshef/SmartHomeAdapters/infoserver

require (
	github.com/golang/protobuf v1.3.1
	github.com/lib/pq v1.0.0
	github.com/mrbenshef/SmartHomeAdapters/microservice v0.0.0
	golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2 // indirect
	google.golang.org/genproto v0.0.0-20190201180003-4b09977fb922 // indirect
	google.golang.org/grpc v1.19.0
)

replace github.com/mrbenshef/SmartHomeAdapters/microservice => ../microservice
