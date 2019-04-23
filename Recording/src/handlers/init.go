/*
@Time : 2019/4/17 16:53 
@Author : kenny zhu
@File : init.go
@Software: GoLand
@Others:
*/
package handlers

import (
	"github.com/kennyzhu/RTPGateway/Recording/src/handlers"
	"fmt"
)

// All handlers init here.
func Init() {
	// default :
	// micro health go.micro.api.gin call this function.
	Handlers.Router.POST("/", NoModules)
	Handlers.Router.GET("/", NoModules)


	// all handlers register here.
	v1 := Handlers.Router.Group("/v1")
	{
		v1.POST("/endpoint/bindleft", handlers.EndpointBindLeftHandler)
		v1.POST("/endpoint/bindright", handlers.EndpointBindRightHandler)
		v1.POST("/endpoint/update", handlers.EndpointUpdateHandler)
		v1.POST("/endpoint/200", handlers.Endpoint200Handler)
		v1.POST("/endpoint/stop", handlers.EndpointStopHandler)


		v1.GET("/endpoint/statis", handlers.EndpointStatisHandler)
		v1.GET("/endpoint/preview", handlers.EndpointPreviewHandler)
	}

	// register other handlers here, each request run in goroutine.
	// To register others...
	fmt.Println("Modules init finished.")
}