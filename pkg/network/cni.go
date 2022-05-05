package network

import (
	"task-start/pkg/ds"

	"github.com/containernetworking/cni/libcni"
)

type CNI struct {
	ds     ds.Datastore
	config *Config
}

func (cni *CNI) GetCniConfig(networkName string) (*libcni.CNIConfig, *libcni.NetworkConfigList, error) {
	cfg, err := cni.ds.GetNetwork(networkName)
	if err != nil {
		return nil, nil, err
	}

	netConf, err := libcni.ConfListFromBytes(cfg.Bytes)
	if err != nil {
		return nil, nil, err
	}

	cninet := &libcni.CNIConfig{
		Path: []string{cni.config.CNIPath},
	}

	return cninet, netConf, nil
}
