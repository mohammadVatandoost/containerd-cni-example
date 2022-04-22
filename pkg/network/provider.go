package network

import (
	"task-start/pkg/ds/local"
)

func NewCNI(conf *Config) (*CNI, error) {
	// "./cniconfigs"
	ds, err := local.New(conf.DsURI)
	if err != nil {
		return nil, err
	}
	return &CNI{
		config: conf,
		ds: ds,
	}, nil
}
