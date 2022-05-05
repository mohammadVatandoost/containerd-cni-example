package main

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/sirupsen/logrus"
	"log"
	"syscall"
	"task-start/internal/app/controller"
	"task-start/pkg/cninetwork"
	cntext "task-start/pkg/context"
)

const (
	NameSpace            = "example"
	RedisContainerID     = "redis-server"
	RedisImageName       = "docker.io/library/redis:alpine"
	MysqlContainerID     = "Mysql-server"
	MysqlImageName       = "mysql:5.7"
	WordpressContainerID = "Wordpress-server"
	WordpressImageName   = "wordpress:beta-php8.1-fpm-alpine"
	NetworkName          = "example"
	CNIPath              = "/opt/cni/bin"
	DSPath               = "./cniconfigs"
	ContainerAddress     = "/run/containerd/containerd.sock"
)

func main() {
	client, err := containerd.New(ContainerAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	cni, err := cninetwork.InitNetwork()
	if err != nil {
		log.Fatal(err)
	}

	controlPlane := controller.NewControlPlane(client, cni, NameSpace)

	err = controlPlane.StartService(RedisImageName, RedisContainerID)
	if err != nil {
		log.Fatal(err)
	}

	serverContext, serverCancel := cntext.WithSignalCancellation(
		context.Background(),
		syscall.SIGTERM, syscall.SIGINT,
	)
	defer serverCancel()
	<-serverContext.Done()

	err = controlPlane.StopService(RedisContainerID)
	if err != nil {
		logrus.Errorf("can not start service, ID: %v, err: %v \n", RedisContainerID, err.Error())
	}

}

