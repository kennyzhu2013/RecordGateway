/*
@Time : 2018/8/21 11:04 
@Author : kenny zhu
@File : init.go
@Software: GoLand
@Others:
*/
package modules

import (
	"router"

	"github.com/kennyzhu/go-os/log"
	"github.com/micro/go-micro"
)

// All handlers init here.
func Init() {
	// default :
	// micro health go.micro.api.gin call this function.
	Modules.Router.POST("/", NoModules)
	Modules.Router.GET("/", NoModules)

	// Media-proxy init client here.
	e := new( mediaProxy )
	wrapper := router.NewClientWrapper("X-Media-Server")
	service := micro.NewService(
		micro.WrapClient(wrapper),
	)
	/*
	serviceWeb :=  web.NewService(
		web.MicroService(serviceBase),
	)*/
	service.Init()
	// Use the generated client stub
	e.cl = service.Client()
	Modules.Router.GET("/Preferences/*action", e.Proxy)
	// Examples End

	// register other handlers here, each request run in goroutine.
	// To register others...

	log.Info("Modules init finished.")
}
