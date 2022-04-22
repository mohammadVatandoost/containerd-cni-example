package service

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

const (
	ContainerAddress = "/run/containerd/containerd.sock"
)

type Service struct {
	Namespace string
	ContainerID string
	Image string
	container containerd.Container
	task containerd.Task
	exitStatusC <- chan containerd.ExitStatus
}

func (s *Service) Run() error {
	client, err := containerd.New(ContainerAddress)
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), s.Namespace)

	image, err := client.Pull(ctx, s.Image, containerd.WithPullUnpack)
	if err != nil {
		return err
	}

	container, err := client.NewContainer(
		ctx,
		s.ContainerID,
		containerd.WithImage(image),
		containerd.WithNewSnapshot(s.ContainerID+"-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}
	s.container = container
	return nil
}

func (s *Service) Close() error {
	ctx := namespaces.WithNamespace(context.Background(), s.Namespace)
	return s.container.Delete(ctx, containerd.WithSnapshotCleanup)
}


