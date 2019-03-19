module github.com/mrbenshef/SmartHomeAdapters/userserver

require (
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/cenkalti/backoff v2.1.1+incompatible // indirect
	github.com/containerd/continuity v0.0.0-20181203112020-004b46473808 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.3.3 // indirect
	github.com/golang/protobuf v1.2.0
	github.com/lib/pq v1.0.0
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/ory/dockertest v3.3.4+incompatible
	github.com/pkg/errors v0.8.1 // indirect
	github.com/sirupsen/logrus v1.3.0 // indirect
	golang.org/x/crypto v0.0.0-20190123085648-057139ce5d2b
	golang.org/x/net v0.0.0-20190206173232-65e2d4e15006
	google.golang.org/grpc v1.18.0
	github.com/mrbenshef/SmartHomeAdapters/microservice v0.0.0
)

replace github.com/mrbenshef/SmartHomeAdapters/microservice => ../microservice
