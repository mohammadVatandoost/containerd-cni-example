package service

import (
	"context"
	"fmt"

	"net"
	api "task-start/api/v1"
	"task-start/pkg/network"

	"github.com/containernetworking/cni/libcni"
	"github.com/pkg/errors"
)

func (s *Service) AddNetwork(networkName string, cninet *libcni.CNIConfig, nc *libcni.NetworkConfigList) (net.IP, error) {
	ctx := context.Background()
	pids, err := s.task.Pids(ctx)
	if err != nil {
		return nil, err
	}

	if len(pids) == 0 {
		return nil, errors.Errorf("no pids found for task in container %s", s.container.ID())
	}

	containerPid := int(pids[0].Pid)
	ifaceName, err := network.GenerateIfaceName(containerPid)
	if err != nil {
		return nil, err
	}

	rt := &libcni.RuntimeConf{
		ContainerID: fmt.Sprintf("%d", containerPid),
		NetNS:       fmt.Sprintf("/proc/%d/ns/net", containerPid),
		IfName:      ifaceName,
	}

	r, err = cninet.AddNetworkList(ctx, nc, rt)
	if err != nil {
		return nil, errors.Wrap(err, "error adding cni network")
	}

	// res, err := current.GetResult(r)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "error getting result from cninet")
	}

	result, err := r.GetAsVersion("0.3.0")
	if err != nil {
		return nil, errors.Wrap(err, "error getting result as version")
	}

	cr := result.(*current.Result)
	if len(cr.IPs) == 0 {
		return nil, fmt.Errorf("container did not receive an IP")
	}

	ipConfig := cr.IPs[0]
	ip := ipConfig.Address.IP

	networkConfig, err := network.LoadNetworkConfig(ctx, s.container)
	if err != nil {
		if err != network.ErrNetworkConfigExtensionNotFound {
			return nil, errors.Wrap(err, "error loading network config")
		}
		networkConfig = &api.NetworkConfig{
			Networks: map[string]*api.ContainerNetworkConfig{},
		}
	}
	if networkConfig.Networks == nil {
		networkConfig.Networks = map[string]*api.ContainerNetworkConfig{}
	}

	networkConfig.Networks[networkName] = &api.ContainerNetworkConfig{
		Interface: ifaceName,
		IP:        ip.String(), 
	}

	if err := container.Update(ctx, network.WithUpdateExtension(network.NetworkConfigExtension, networkConfig)); err != nil {
		return nil, errors.Wrap(err, "error updating container extension")
	}
	
	return ip, nil
}
