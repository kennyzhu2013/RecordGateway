/*
@Time : 2018/8/21 11:04 
@Author : kenny zhu
@File : init.go
@Software: GoLand
@Others:
*/
package modules

import (
	"github.com/micro/go-micro/client"
	example "github.com/kennyzhu/go-os/dbservice/proto/example"
	"conf"

	"github.com/kennyzhu/go-os/log"
	"github.com/kennyzhu/go-os/plugins/etcdv3"
)

// All handlers init here.
func Init() {
	// default :
	// micro health go.micro.api.gin call this function.
	Modules.Router.POST("/", NoModules)
	Modules.Router.GET("/", NoModules)

	// Examples Begin:micro api --handler=http as proxy, default is rpc .
	// Base module router for rest full, Preferences is name of table or tables while Module equals database.
	// Call url:curl "http://localhost:8004/Preferences/GetPreference?user=1"
	e := new( examples )
	client.DefaultClient = client.NewClient( client.Registry(etcdv3.DefaultEtcdRegistry) )
	e.cl = example.NewPreferencesService(conf.ApiConf.SrvName, client.DefaultClient)
	Modules.Router.GET("/Preferences/*action", e.Preferences)
	// Examples End

	// register other handlers here, each request run in goroutine.
	// To register others...

	log.Info("Modules init finished.")
}
