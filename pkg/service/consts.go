package service

const (
	ContainerAddress = "/run/containerd/containerd.sock"
	SnapShotter      = "overlayfs"
	HostDir          = "/var/lib/myaws"
	MountsDir        = HostDir + "/mounts"
)
