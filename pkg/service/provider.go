package service

import "github.com/containerd/containerd"

func New(containerdClient *containerd.Client, namespace string,
	containerID string, image string) *Service {
	return &Service{
		Namespace:        namespace,
		ContainerID:      containerID,
		Image:            image,
		containerdClient: containerdClient,
	}
}
