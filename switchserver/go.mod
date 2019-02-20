module github.com/mrbenshef/SmartHomeAdapters/switchserver

require (
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/cenkalti/backoff v2.1.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20181203112020-004b46473808 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.3.3 // indirect
	github.com/golang/protobuf v1.2.0
	github.com/lib/pq v1.0.0
	github.com/mrbenshef/SmartHomeAdapters/robotserver v0.0.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/ory/dockertest v3.3.4+incompatible
	github.com/pkg/errors v0.8.1 // indirect
	github.com/sirupsen/logrus v1.3.0 // indirect
	golang.org/x/net v0.0.0-20180826012351-8a410e7b638d
	golang.org/x/sys v0.0.0-20190201152629-afcc84fd7533 // indirect
	google.golang.org/grpc v1.18.0
	gopkg.in/h2non/gock.v1 v1.0.14 // indirect
)

replace github.com/mrbenshef/SmartHomeAdapters/robotserver => ../robotserver
