package controller

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"task-start/pkg/service"
	"time"
)

func (cp *ControlPlane) StartService(ID string, name string) error {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	_, ok := cp.services[name]
	if ok {
		return fmt.Errorf("can not start the service, because service with name %v is exist", name)
	}

	start := time.Now()
	testService := service.New(cp.containerdClient, cp.NameSpace, name, ID)
	err := testService.Run()
	if err != nil {
		return err
	}

	err = testService.StartTask(cp.cni)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("cold start: %v \n", time.Now().Sub(start).Milliseconds())

	return nil
}

func (cp *ControlPlane) StopService(name string) error {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	Service, ok := cp.services[name]
	if ok {
		return fmt.Errorf("can not stop the service, because service with name %v is not exist", name)
	}

	err := Service.StopTask()
	if err != nil {
		logrus.Warnf("can not stop task, name: %v, err: %v", name, err.Error())
	}

	err = Service.Close()
	if err != nil {
		logrus.Warnf("can not stop service, name: %v, err: %v", name, err.Error())
		return err
	}
	delete(cp.services, name)
	return nil
}
