package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	. "modules"
	"conf"

	"github.com/micro/go-micro/registry"
)

func main() {
	initService()

	// created in goroutine by gin.
	go Modules.Router.Run( conf.AppConf.Address )

	// Register modules and app.Run...
	// All path processed by modules..
	// service.Handle("/", Modules.App)
	registry.Register(service)

	// 通过registry可以获得服务器的ip和端口等信息...
	// find self
	rsp, err := registry.GetService(service.Name)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Got service %+v\n", rsp[0])
		fmt.Printf("Nodes info %+v\n", rsp[0].Nodes[0])
	}

	/*
	service.Name = "go.micro.api.gin-gateway2"
    registry.Register(service)
	rsp, err = registry.GetService(service.Name)
    if err != nil {
        fmt.Println(err)
    } else {
        fmt.Printf("Got service2 %+v\n", rsp[0])
        fmt.Printf("Nodes info2 %+v\n", rsp[0].Nodes[0])
    }*/

	// micro health查询需要export MICRO_PROXY_ADDRESS=0.0.0.0:8002支持http json方式访问..
	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-notify

	registry.Deregister(service)
}
