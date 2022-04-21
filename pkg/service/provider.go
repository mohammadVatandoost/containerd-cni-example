package service

func New(namespace string, containerID string, image string) *Service {
	return &Service{
		Namespace: namespace,
		ContainerID: containerID,
		Image: image,
	}
}
