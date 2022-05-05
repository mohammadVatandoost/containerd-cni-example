package main

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/sirupsen/logrus"
	"log"
	"path"
	"syscall"
	"task-start/internal/app/controller"
	"task-start/pkg/cninetwork"
	cntext "task-start/pkg/context"
	"task-start/pkg/service"
)

const (
	NameSpace            = "example"
	RedisContainerID     = "redis-server"
	RedisImageName       = "docker.io/library/redis:alpine"
	MysqlContainerID     = "Mysql-server"
	MysqlImageName       = "docker.io/library/mysql:5.7"
	WordpressContainerID = "Wordpress-server"
	WordpressImageName   = "docker.io/library/wordpress:beta-php8.1-fpm-alpine"
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

	err = controlPlane.StartService(MysqlImageName, MysqlContainerID, getMYSQLEnv(), getMYSQLMounts())
	if err != nil {
		log.Fatal(err)
	}

	serverContext, serverCancel := cntext.WithSignalCancellation(
		context.Background(),
		syscall.SIGTERM, syscall.SIGINT,
	)
	defer serverCancel()
	<-serverContext.Done()

	err = controlPlane.StopService(MysqlContainerID)
	if err != nil {
		logrus.Errorf("can not start service, ID: %v, err: %v \n", MysqlContainerID, err.Error())
	}

}

func getMYSQLEnv() map[string]string {
	env := make(map[string]string)
	env["MYSQL_ROOT_PASSWORD"] = "somewordpress"
	env["MYSQL_DATABASE"] = "wordpress"
	env["MYSQL_USER"] = "wordpress"
	env["MYSQL_PASSWORD"] = "wordpress"
	return env
}

func getMYSQLMounts() []specs.Mount {
	var mounts []specs.Mount
	mounts = append(mounts, specs.Mount{
		Destination: path.Join(service.MountsDir, "db_data"),
		Type:        "bind",
		Source:      "/var/lib/mysql",
		Options:     []string{"rwbind", "rw"},
	})
	return mounts
}
