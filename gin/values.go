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
)

// registry service ip and port to the ET-CD.
// address and port must be re-written.
var (
	service = &registry.Service{
		Name: "go.micro.api.gin-gateway",
		Nodes: []*registry.Node{
			{
				Id:      "go.micro.api.gin-gateway-" + uuid.NewUUID().String(),
				Address: "localhost",
				Port:    8400,
			},
		},
	}
)