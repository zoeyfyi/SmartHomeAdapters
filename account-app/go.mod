module github.com/mrbenshef/SmartHomeAdapters/account-app

require (
        github.com/golang/protobuf v1.2.1-0.20181127190454-8d0c54c12466
        github.com/julienschmidt/httprouter v1.2.0
        github.com/mrbenshef/SmartHomeAdapters/userserver v0.0.0
        google.golang.org/grpc v1.18.0
)


replace github.com/mrbenshef/SmartHomeAdapters/userserver => ../userserver

