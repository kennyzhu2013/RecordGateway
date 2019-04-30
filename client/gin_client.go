/*
@Time : 2019/4/17 15:38
@Author : kenny zhu
@File : gin_client.go
@Software: GoLand
@Others:
*/
package main

import (
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"

	example "github.com/kennyzhu/go-os/dbservice/proto/example"
	"github.com/micro/go-web"
	"github.com/kennyzhu/go-os/plugins/etcdv3"
	"github.com/micro/go-micro/registry"
)

func main()  {
	// r, err := http.Get("http://localhost:8002/dbservice/Preferences/PreferenceList?limit=2&index=1")
	registry.DefaultRegistry = etcdv3.NewRegistry()
	service := web.NewService()
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	c := service.Client()

	// not need micro web, call by registry
	// can be called by any other server, eg: proxy go.micro.api.gin.
	r,err := c.Get("http://go.micro.api.media-gateway/Preferences/GetPreferencesList?limit=2&index=1")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(b)
	fmt.Println(r.StatusCode)
	var body map[string]interface{}
	if err := json.Unmarshal(b, &body); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(body["message"])
	rsp := make([]example.Preference, 0)
	if err := json.Unmarshal( []byte( body["message"].(string) ) , &rsp); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp)
}
