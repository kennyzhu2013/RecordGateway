package main

import (
	"config"
	"core"
	"github.com/gin-gonic/gin"
	//"gopkg.in/yaml.v2"
)
import "handlers"

var Fuck int

func main() {
	// init config
	config.InitConfig("setting.yaml")

	//init the session manager
	core.InitSessionManage(config.GetConfig().Ims.Ip)

	//start http api server
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/endpoint/bindleft", handlers.EndpointBindLeftHandler)
		v1.POST("/endpoint/bindright", handlers.EndpointBindRightHandler)
		v1.POST("/endpoint/update", handlers.EndpointUpdateHandler)
		v1.POST("/endpoint/200", handlers.Endpoint200Handler)
		v1.POST("/endpoint/stop", handlers.EndpointStopHandler)


		v1.GET("/endpoint/statis", handlers.EndpointStatisHandler)
		v1.GET("/endpoint/preview", handlers.EndpointPreviewHandler)
	}

	router.Run(config.GetConfig().Http.Address)
}
