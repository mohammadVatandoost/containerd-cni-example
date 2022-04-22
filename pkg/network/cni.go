package network

import (
	"github.com/containernetworking/cni/libcni"
	"task-start/pkg/ds"
)


type CNI struct {
	ds         ds.Datastore
	config     *Config
}

func (cni *CNI) getCniConfig(networkName string) (*libcni.CNIConfig, *libcni.NetworkConfigList, error) {
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
