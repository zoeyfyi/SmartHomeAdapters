module github.com/mrbenshef/SmartHomeAdapters/switchserver

require (
	github.com/golang/protobuf v1.2.0
	github.com/lib/pq v1.0.0
	github.com/mrbenshef/SmartHomeAdapters/robotserver v0.0.0
	golang.org/x/net v0.0.0-20180826012351-8a410e7b638d
	golang.org/x/sys v0.0.0-20190201152629-afcc84fd7533 // indirect
	google.golang.org/grpc v1.18.0
	gopkg.in/h2non/gock.v1 v1.0.14
)

replace github.com/mrbenshef/SmartHomeAdapters/robotserver => ../robotserver
