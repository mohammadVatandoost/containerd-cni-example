module task-start

go 1.16

require github.com/containerd/containerd v1.6.2

require (
	github.com/containerd/go-cni v1.1.3
	github.com/containerd/typeurl v1.0.2
	github.com/containernetworking/cni v1.0.1
	github.com/gogo/protobuf v1.3.2
	github.com/google/martian v2.1.0+incompatible
	github.com/opencontainers/runtime-spec v1.0.3-0.20210326190908-1c3f411f0417
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/vishvananda/netlink v1.1.1-0.20210330154013-f5de75959ad5
	github.com/vishvananda/netns v0.0.0-20210104183010-2eb08e3e575f
	google.golang.org/grpc v1.43.0
)
