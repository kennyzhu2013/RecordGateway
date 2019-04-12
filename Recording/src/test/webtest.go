package main

import "github.com/gin-gonic/gin"
import "handlers"

func main() {
	router := gin.Default()

	v1 := router.Group("/v1")
	{
		v1.POST("/endpoint/apply", handlers.EndpointBindLeftHandler) //allocate
		v1.POST("/endpoint/start", handlers.EndpointBindRightHandler)
		v1.GET("/endpoint/preview", handlers.EndpointPreviewHandler)
	}

	router.Run(":8080")
}
