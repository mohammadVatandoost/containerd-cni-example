package network

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/containers"
	"github.com/containerd/containerd/errdefs"
	"github.com/containerd/typeurl"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
	"strings"
	api "task-start/api/v1"
)

const (
	maxInterfaceCount = 10
	// NetworkConfigExtension is the name of the containerd extension to store the config
	NetworkConfigExtension = "io.circuit.network"
	// RestartLabel is the label used for automatically restarting a container upon stop
	RestartLabel = "io.circuit.restart"
)

var (
	// ErrNetworkConfigExtensionNotFound is returned when the network config containerd extension is not found
	ErrNetworkConfigExtensionNotFound = errors.New("network config extension not found")
)

func GenerateIfaceName(containerPid int) (string, error) {
	originalNs, err := netns.Get()
	if err != nil {
		return "", err

	}
	defer originalNs.Close()

	cntNs, err := netns.GetFromPid(containerPid)
	if err != nil {
		return "", err
	}
	defer cntNs.Close()

	ifaceName := ""
	netns.Set(cntNs)
	for i := 0; i < maxInterfaceCount; i++ {
		n := fmt.Sprintf("eth%d", i)
		if _, err := netlink.LinkByName(n); err != nil {
			if !strings.Contains(err.Error(), "no such network interface") {
				ifaceName = n
				break
			}
		}
	}
	netns.Set(originalNs)

	if ifaceName == "" {
		return "", fmt.Errorf("unable to generate device name; maximum number of devices reached (%d)", maxInterfaceCount)
	}

	return ifaceName, nil
}


func LoadNetworkConfig(ctx context.Context, c containerd.Container) (*api.NetworkConfig, error) {
	extensions, err := c.Extensions(ctx)
	if err != nil {
		return nil, err
	}

	ext, ok := extensions[NetworkConfigExtension]
	if !ok {
		return nil, ErrNetworkConfigExtensionNotFound
	}

	v, err := typeurl.UnmarshalAny(&ext)
	if err != nil {
		return nil, err
	}

	e, ok := v.(*api.NetworkConfig)
	if !ok {
		return nil, errors.Errorf("expected type 'v1.NetworkConfig'; received %T", v)
	}

	return e, nil
}

func WithUpdateExtension(name string, extension interface{}) containerd.UpdateContainerOpts {
	return func(ctx context.Context, _ *containerd.Client, c *containers.Container) error {
		if name == "" {
			return errors.Wrapf(errdefs.ErrInvalidArgument, "extension name must not be zero-length")
		}
		any, err := typeurl.MarshalAny(extension)
		if err != nil {
			if errors.Cause(err) == typeurl.ErrNotFound {
				return errors.Wrapf(err, "extension %q is not registered with the typeurl package, see `typeurl.Register`", name)
			}
			return errors.Wrap(err, "error marshalling extension")
		}

		if c.Extensions == nil {
			c.Extensions = make(map[string]types.Any)
		}
		c.Extensions[name] = *any
		return nil
	}
}
