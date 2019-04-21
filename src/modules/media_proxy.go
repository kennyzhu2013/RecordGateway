/*
@Time : 2018/8/21 11:19 
@Author : kenny zhu
@File : proxy.go
@Software: GoLand
@Others:
*/
package modules

import (
	"github.com/kennyzhu/go-os/log"
	// proto "proto"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/client"
	"conf"
	"io/ioutil"

	"github.com/micro/go-micro/metadata"
)

type mediaProxy struct{
	cl client.Client
}

// All are run in goroutine
func (s *mediaProxy) Proxy(ctx *gin.Context) {
	action := ctx.Param("action")

	switch action {
	case "/Invite":
		s.reverseProxy( ctx, action )
	default:
		/*
		ctx.JSON(404, map[string]string {
			"message": "Unknown action:" + action,
		})*/
		s.reverseProxy( ctx, action )
	}
	log.Debug("Proxy done!")
}

func (s *mediaProxy) reverseProxy(ctx *gin.Context, action string) {
	// todo: 获取请求所有头部并全部写入到context中去..
	// push map[string]string.
	headers := make(map[string]string, 20)
	for k,v :=range ctx.Request.Header {
		if "X-Media-Server" == k {
			for _,address := range v {
				headers[k] += address
			}
			break
		}
	}

	// get request url..
	b, _ := ioutil.ReadAll(ctx.Request.Body)
	url := conf.ApiConf.SrvName + string("/Preferences") + action
	log.Infof("Received reverseProxy http request:%v", method)

	// push http url
	headers["reverseProxy_url"] = url
	var headersContext context.Context = metadata.NewContext(context.Background(), headers )

	// modify request body ....
	// NewRequest
	req := s.cl.NewRequest(conf.ApiConf.SrvName, method, string(b), client.WithContentType("application/json"))
	var rsp string

	// Call service
	if err := client.Call(headersContext, req, &rsp); err != nil {
		ctx.JSON(500, map[string]string{
			"message": err.Error(),
		})
		log.Error(err)
		return
	}

	ctx.JSON(200, rsp)
	log.Info("reverseProxy End:")
}
