package controller

import (
	"github.com/containerd/containerd"
	gocni "github.com/containerd/go-cni"
	"sync"
	"task-start/pkg/service"
)

type ControlPlane struct {
	containerdClient *containerd.Client
	services         map[string]*service.Service
	cni              gocni.CNI
	NameSpace        string
	lock             sync.Mutex
}
