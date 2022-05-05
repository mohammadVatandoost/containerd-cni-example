package main

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"task-start/pkg/cninetwork"
	"task-start/pkg/service"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

const (
	NameSpace        = "example"
	ContainerID      = "redis-server"
	ImageName        = "docker.io/library/redis:alpine"
	NetworkName      = "example"
	CNIPath          = "/opt/cni/bin"
	DSPath           = "./cniconfigs"
	ContainerAddress = "/run/containerd/containerd.sock"
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

	testService := service.New(client, NameSpace, ContainerID, ImageName)
	err = testService.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer testService.Close()

	err = testService.StartTask(cni)
	if err != nil {
		log.Fatal(err)
	}

	//cniNetwork, err := network.NewCNI(&network.Config{
	//	CNIPath: CNIPath,
	//	DsURI:   DSPath,
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	//cninet, nc, err := cniNetwork.GetCniConfig(NetworkName)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//ip, err := testService.AddNetwork(NetworkName, cninet, nc)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("container IP: %v \n", ip.String())

	time.Sleep(20 * time.Second)
	testService.StopTask()

	// if err := redisExample(); err != nil {
	// 	log.Fatal(err)
	// }
}

func redisExample() error {
	// create a new client connected to the default socket path for containerd
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	// create a new context with an "example" namespace
	ctx := namespaces.WithNamespace(context.Background(), "example")

	// pull the redis image from DockerHub
	image, err := client.Pull(ctx, "docker.io/library/redis:alpine", containerd.WithPullUnpack)
	if err != nil {
		return err
	}

	// create a container
	container, err := client.NewContainer(
		ctx,
		"redis-server",
		containerd.WithImage(image),
		containerd.WithNewSnapshot("redis-server-snapshot", image),
		containerd.WithNewSpec(oci.WithImageConfig(image)),
	)
	if err != nil {
		return err
	}
	defer container.Delete(ctx, containerd.WithSnapshotCleanup)

	// create a task from the container
	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	defer task.Delete(ctx)
	time.Sleep(2 * time.Second)
	fmt.Println("******* task.Wait ********")
	// make sure we wait before calling start
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("******* task.Start ********")
	// call start on the task to execute the redis server
	if err := task.Start(ctx); err != nil {
		return err
	}
	fmt.Println("============")
	// sleep for a lil bit to see the logs
	time.Sleep(3 * time.Second)

	fmt.Println("***************")
	time.Sleep(3 * time.Second)
	// kill the process and get the exit status
	if err := task.Kill(ctx, syscall.SIGTERM); err != nil {
		return err
	}

	// wait for the process to fully exit and print out the exit status

	status := <-exitStatusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}
	fmt.Printf("redis-server exited with status: %d\n", code)

	return nil
}
