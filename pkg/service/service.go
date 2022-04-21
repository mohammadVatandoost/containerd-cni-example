package service

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"syscall"
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

func (s *Service) StartTask() error {
	ctx := namespaces.WithNamespace(context.Background(), s.Namespace)
	task, err := s.container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	s.task = task

	exitStatus, err := task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// call start on the task to execute the redis server
	if err := task.Start(ctx); err != nil {
		return err
	}
	s.exitStatusC = exitStatus
	return nil
}

func (s *Service) StopTask() error {
	ctx := namespaces.WithNamespace(context.Background(), s.Namespace)
	if err := s.task.Kill(ctx, syscall.SIGTERM); err != nil {
		return err
	}

	status := <-s.exitStatusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	fmt.Printf("container : %v exited with status: %d \n", s.ContainerID, code)
	exitStatus, err := s.task.Delete(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("container : %v, task exitStatus: %v \n", s.ContainerID, exitStatus)
	return nil
}
