package service

import (
	"github.com/containerd/containerd"
	"github.com/opencontainers/runtime-spec/specs-go"
)

func New(containerdClient *containerd.Client, namespace string,
	containerID string, image string, envVars map[string]string,
	volumeMounts []specs.Mount) *Service {
	return &Service{
		Namespace:        namespace,
		ContainerID:      containerID,
		Image:            image,
		containerdClient: containerdClient,
		reqEnvVars:       envVars,
		volumeMounts:     volumeMounts,
	}
}
