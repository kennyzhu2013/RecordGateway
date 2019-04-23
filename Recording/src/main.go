package main

import (
	. "config"
	"core"
	//  "gopkg.in/yaml.v2"
)
import (
	. "handlers"
)

var Fuck int

func main() {
	initService()

	// init the session manager
	core.InitSessionManage(AppConf.Ims.Ip)

	// start http api server
	// router := gin.Default()
	go Handlers.Router.Run(AppConf.HttpAddress)

	registryStart()
}
