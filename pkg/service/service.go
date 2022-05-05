package service

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/google/martian/log"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type Service struct {
	Namespace        string
	ContainerID      string
	Image            string
	container        containerd.Container
	task             containerd.Task
	containerdClient *containerd.Client
	exitStatusC      <-chan containerd.ExitStatus
	reqEnvVars       map[string]string
	volumeMounts     []specs.Mount
}

func (s *Service) Run() error {

	ctx := namespaces.WithNamespace(context.Background(), s.Namespace)

	image, err := GetImage(ctx, s.containerdClient, s.Image)
	if err != nil {
		return err
	}

	containerInfo, err := s.containerdClient.ContainerService().Get(ctx, s.ContainerID)
	if err == nil {
		log.Infof("container exist, ContainerID: %v", containerInfo.ID)
		s.container, err = s.containerdClient.LoadContainer(ctx, containerInfo.ID)
		if err != nil {
			log.Errorf("can not load container, ContainerID: %v", containerInfo.ID)
			return err
		}
		return nil
	}
	envs := prepareEnv(s.reqEnvVars)

	mounts := getOSMounts()
	mounts = append(mounts, s.volumeMounts...)
	//for _, secret := range req.Secrets {
	//	mounts = append(mounts, specs.Mount{
	//		Destination: path.Join("/var/openfaas/secrets", secret),
	//		Type:        "bind",
	//		Source:      path.Join(secretMountPath, secret),
	//		Options:     []string{"rbind", "ro"},
	//	})
	//}

	log.Infof("Creating container, ContainerID: %v", containerInfo.ID)
	//snapShotter := client.SnapshotService(SnapShotter)
	//ToDo: handle snapshotter is exist
	container, err := s.containerdClient.NewContainer(
		ctx,
		s.ContainerID,
		containerd.WithImage(image),
		//containerd.WithSnapshot(s.ContainerID+"-snapshot"),
		containerd.WithNewSnapshot(s.ContainerID+"-snapshot", image),
		containerd.WithNewSpec(
			oci.WithHostname(s.ContainerID),
			oci.WithImageConfig(image),
			oci.WithEnv(envs),
			oci.WithMounts(mounts),
		),
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
