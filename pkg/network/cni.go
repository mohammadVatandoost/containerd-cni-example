package network

import (
	"fmt"
	"github.com/containernetworking/cni/libcni"
	"task-start/pkg/ds"
)


type CNI struct {
	ds         ds.Datastore
	config     *Config
}

func (cni *CNI) getCniConfig(networkName string, containerPid int, ifaceName string) (*libcni.CNIConfig, *libcni.NetworkConfigList, *libcni.RuntimeConf, error) {
	cfg, err := cni.ds.GetNetwork(networkName)
	if err != nil {
		return nil, nil, nil, err
	}

	netConf, err := libcni.ConfListFromBytes(cfg.Bytes)
	if err != nil {
		return nil, nil, nil, err
	}

	cninet := &libcni.CNIConfig{
		Path: []string{cni.config.CNIPath},
	}

	rt := &libcni.RuntimeConf{
		ContainerID: fmt.Sprintf("%d", containerPid),
		NetNS:       fmt.Sprintf("/proc/%d/ns/net", containerPid),
		IfName:      ifaceName,
	}

	return cninet, netConf, rt, nil
}
