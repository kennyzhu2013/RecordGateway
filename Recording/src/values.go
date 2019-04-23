/*
@Time : 2019/4/17 16:40 
@Author : kenny zhu
@File : values.go
@Software: GoLand
@Others:
*/
package main

import (
	"github.com/pborman/uuid"
	. "config"
	. "modules"
	"github.com/micro/go-micro/registry"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// registry service ip and port to the ET-CD.
// address and port must be re-written.
var (
	service = &registry.Service{
		Name: "go.micro.api.media-proxy",
		Metadata: map[string]string{
			"serverDescription": "audio recording proxy service",  // server desc.
		},
		Nodes: []*registry.Node{
			{
				Id:      "go.micro.api.media-proxy-" + uuid.NewUUID().String(),
				Address: "localhost",
				Port:    8400,
				Metadata: map[string]string{
					"serverTag": "media-proxy",  // server division.
				},
			},
		},
	}
)

func initService()  {
	// init config
	InitConfig("setting.yaml")

	service.Name = AppConf.Api.SrvName
	var nodeSelf = service.Nodes[0]
	nodeSelf.Address = AppConf.Http.Ip
	nodeSelf.Port = AppConf.Http.Port

	Init()
}

func registryStart()  {
	// Register modules and app.Run...
	// Here default registry is consul.
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

	// micro health查询需要export MICRO_PROXY_ADDRESS=0.0.0.0:8002支持http json方式访问..
	notify := make(chan os.Signal, 1)
	signal.Notify(notify, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-notify

	registry.Deregister(service)
}