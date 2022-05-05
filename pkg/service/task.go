package service

import (
	"context"
	"fmt"
	"log"

	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	gocni "github.com/containerd/go-cni"
	"syscall"
	"task-start/pkg/cninetwork"
)

func (s *Service) StartTask(cni gocni.CNI) error {
	ctx := namespaces.WithNamespace(context.Background(), s.Namespace)
	//ToDo handle remove extera load container

	task, err := s.container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	s.task = task

	name := s.container.ID()

	log.Printf("Container ID: %s\tTask ID %s:\tTask PID: %d\t\n", name, task.ID(), task.Pid())

	labels := map[string]string{}
	_, err = cninetwork.CreateCNINetwork(ctx, cni, task, labels)

	if err != nil {
		return err
	}

	ip, err := cninetwork.GetIPAddress(name, task.Pid())
	if err != nil {
		return err
	}

	log.Printf("%s has IP: %s.\n", name, ip)

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
