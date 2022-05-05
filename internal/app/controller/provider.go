package controller

import (
	"github.com/containerd/containerd"
	gocni "github.com/containerd/go-cni"
	"task-start/pkg/service"
)

func NewControlPlane(containerdClient *containerd.Client, cni gocni.CNI, NameSpace string) *ControlPlane {
	return &ControlPlane{
		containerdClient: containerdClient,
		services:         make(map[string]*service.Service),
		cni:              cni,
		NameSpace:        NameSpace,
	}
}
