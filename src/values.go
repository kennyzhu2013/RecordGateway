/*
@Time : 2019/4/15 16:51 
@Author : kenny zhu
@File : values.go
@Software: GoLand
@Others:
*/
package main

import (
	"github.com/micro/go-micro/registry"
	"github.com/pborman/uuid"
	"conf"
	. "modules"

	"github.com/kennyzhu/go-os/plugins/etcdv3"
)

// registry service ip and port to the ET-CD.
// address and port must be re-written.
var (
	service = &registry.Service{
		Name: "go.micro.api.gin-gateway",
		Metadata: map[string]string{
			"serverTag": "media-gateway",  // service desc.
		},
		Endpoints: []*registry.Endpoint {
			{
				Name: "Preferences.GetPreferencesList",
			},
		},
		Nodes: []*registry.Node{
			{
				Id:      "go.micro.api.gin-gateway-" + uuid.NewUUID().String(),
				Address: "localhost",
				Port:    8400,
				Metadata: map[string]string{
					"serverTag": "media-gateway",  // node division.
				},
			},
		},
	}
)

func initService()  {
	conf.Init( "./conf/gin-api.json" )

	service.Name = conf.ApiConf.ApiName
	var nodeSelf = service.Nodes[0]
	nodeSelf.Address = conf.AppConf.IP
	nodeSelf.Port = conf.AppConf.Port

	registry.DefaultRegistry = etcdv3.DefaultEtcdRegistry
	Init()
}