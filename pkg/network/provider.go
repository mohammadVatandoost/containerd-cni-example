package network

import (
	"task-start/pkg/ds/local"
)

func NewCNI(conf *Config) (*CNI, error) {
	ds, err := local.New("./cniconfigs")
	if err != nil {
		return nil, err
	}
	return &CNI{
		config: conf,
		ds: ds,
	}, nil
}
