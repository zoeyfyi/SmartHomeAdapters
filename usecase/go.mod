module github.com/mrbenshef/SmartHomeAdapters/usecase

require (
	github.com/mrbenshef/SmartHomeAdapters/robotserver v0.0.0
	google.golang.org/grpc v1.18.0
)

replace github.com/mrbenshef/SmartHomeAdapters/robotserver => ../robotserver
