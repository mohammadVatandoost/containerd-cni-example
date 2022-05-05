package service

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/sirupsen/logrus"
)

func GetImage(ctx context.Context, containerdClient *containerd.Client,
	ImageID string) (containerd.Image, error) {
	image, err := containerdClient.GetImage(ctx, ImageID)
	if err == nil {
		return image, nil
	}
	logrus.Infof("Image does not exist, lets pull it, imageID: %v \n", ImageID)
	image, err = containerdClient.Pull(ctx, ImageID, containerd.WithPullUnpack)
	if err != nil {
		return nil, err
	}
	return image, nil
}

//func (s *Service) LoadContainer() error {
//
//}
