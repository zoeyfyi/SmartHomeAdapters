module github.com/mrbenshef/SmartHomeAdapters/clientserver

require (
	github.com/golang/protobuf v1.2.1-0.20181127190454-8d0c54c12466
	github.com/julienschmidt/httprouter v1.2.0
	github.com/mrbenshef/SmartHomeAdapters/infoserver v0.0.0
	github.com/mrbenshef/SmartHomeAdapters/userserver v0.0.0
	google.golang.org/grpc v1.18.0
)

replace github.com/mrbenshef/SmartHomeAdapters/infoserver => ../infoserver

replace github.com/mrbenshef/SmartHomeAdapters/userserver => ../userserver

replace github.com/mrbenshef/SmartHomeAdapters/switchserver => ../switchserver

replace github.com/mrbenshef/SmartHomeAdapters/thermostatserver => ../thermostatserver

replace github.com/mrbenshef/SmartHomeAdapters/robotserver => ../robotserver
