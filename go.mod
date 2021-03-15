module github.com/allgdante/docker-multilogger-plugin

require (
	cloud.google.com/go v0.52.0 // indirect
	cloud.google.com/go/logging v1.0.0 // indirect
	github.com/Graylog2/go-gelf v0.0.0-20170811154226-7ebf4f536d8f // indirect
	github.com/Microsoft/go-winio v0.4.15-0.20190919025122-fc70bd9a86b5 // indirect
	github.com/RackSec/srslog v0.0.0-20180709174129-a4725f04ec91 // indirect
	github.com/allgdante/srslog v0.0.0-20200314183408-5c1512acc434
	github.com/aws/aws-sdk-go v1.28.13 // indirect
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/bsphere/le_go v0.0.0-20200109081728-fc06dab2caa8 // indirect
	github.com/containerd/containerd v1.3.3 // indirect
	github.com/containerd/fifo v0.0.0-20191213151349-ff969a566b00
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e // indirect
	github.com/coreos/go-systemd/v22 v22.2.0 // indirect
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-plugins-helpers v0.0.0-20200102110956-c9a8a2d92ccc
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fluent/fluent-logger-golang v1.5.0 // indirect
	github.com/fsnotify/fsnotify v1.4.7 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/hashicorp/go-multierror v1.0.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/philhofer/fwd v1.0.0 // indirect
	github.com/prometheus/client_golang v1.4.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.5.1
	github.com/tinylib/msgp v1.1.2 // indirect
	go.opencensus.io v0.22.3 // indirect
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
	google.golang.org/api v0.17.0 // indirect
	google.golang.org/genproto v0.0.0-20200207204624-4f3edf09f4f6 // indirect
	google.golang.org/grpc v1.27.1 // indirect
	gotest.tools/v3 v3.0.3 // indirect
)

replace (
	github.com/Graylog2/go-gelf => gopkg.in/Graylog2/go-gelf.v2 v2.0.0-20191017102106-1550ee647df0
	github.com/docker/docker => github.com/moby/moby v20.10.5+incompatible
)

go 1.13
