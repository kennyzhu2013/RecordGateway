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

	"github.com/gin-gonic/gin"
	"conf"
	"io/ioutil"

	. "router"
)

type mediaProxy struct{
	cl WebClient
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
	// todo: 获取请求所有头部并全部写入到服务器上， push map[string]string.
	var destination string
	for k,v :=range ctx.Request.Header {
		if "X-Media-Server" == k {
			for _,address := range v {
				destination += address
			}
			break
		}
	}

	// get request url..
	b, _ := ioutil.ReadAll(ctx.Request.Body)
	serviceUrl := conf.ApiConf.SrvName + string("/Preferences") + action + "?limit=2&index=1"
	log.Infof("Received reverseProxy http request:%v", serviceUrl)

	// Post http request to the destination ....
	// rsp,err := s.cl.Post(serviceUrl, "application/json", destination, bytes.NewReader(b))
	b = b
	rsp,err := s.cl.TestGet(serviceUrl)

	// Call service
	if err != nil {
		ctx.JSON(500, err.Error())
		log.Error(err)
		return
	}
	defer rsp.Body.Close()

	rBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		ctx.JSON(500, err.Error())
		log.Error(err)
		return
	}

	// send body..
	var sBody = string(rBody)
	ctx.Data(rsp.StatusCode, "application/json", rBody)

	log.Infof("reverseProxy End with http body:%v", sBody)
}
