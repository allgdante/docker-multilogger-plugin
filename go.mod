module github.com/allgdante/docker-multilogger-plugin

require (
	cloud.google.com/go v0.86.0 // indirect
	cloud.google.com/go/logging v1.4.2 // indirect
	github.com/Graylog2/go-gelf v0.0.0-20170811154226-7ebf4f536d8f // indirect
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/RackSec/srslog v0.0.0-20180709174129-a4725f04ec91 // indirect
	github.com/aws/aws-sdk-go v1.39.0 // indirect
	github.com/bsphere/le_go v0.0.0-20200109081728-fc06dab2caa8 // indirect
	github.com/containerd/containerd v1.5.2 // indirect
	github.com/containerd/fifo v1.0.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/docker/docker v20.10.7+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-plugins-helpers v0.0.0-20210623094020-7ef169fb8b8e
	github.com/docker/go-units v0.4.0
	github.com/fluent/fluent-logger-golang v1.6.1 // indirect
	github.com/gogo/protobuf v1.3.2
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/prometheus/common v0.29.0 // indirect
	github.com/prometheus/procfs v0.7.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.6.1
	github.com/tinylib/msgp v1.1.6 // indirect
	golang.org/x/net v0.7.0 // indirect
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6 // indirect
	google.golang.org/genproto v0.0.0-20210701191553-46259e63a0a9 // indirect
	google.golang.org/grpc v1.39.0 // indirect
)

replace github.com/Graylog2/go-gelf => gopkg.in/Graylog2/go-gelf.v2 v2.0.0-20191017102106-1550ee647df0

go 1.13
