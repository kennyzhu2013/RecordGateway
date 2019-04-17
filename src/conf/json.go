/*
@Time : 2018/8/21 11:28 
@Author : kenny zhu
@File : json.go
@Software: GoLand
@Others:
*/
package conf

import (
	"github.com/micro/go-config"
	"fmt"
	"github.com/kennyzhu/go-os/log"
	"github.com/micro/go-config/source/file"
	"strconv"
)

var AppConf struct{
	// logger..
	LogLevel int32 `json:"LogLevel"`
	LogPath  string    `json:"LogPath"`

	// app..
	IP string     `json:"IP"`
	Port 	int     `json:"Port"`
	VersionTag int     `json:"Version"`

	Address string
	// Todo:
}

var ApiConf struct{
	// srv set..
	ApiName string     `json:"ApiName"`
	SrvName string     `json:"SrvName"`

	// Todo:
}

// self init...configFile = "./conf/api.json"...
func Init(configFile string) {
	// read config
	if err := config.Load(file.NewSource(
		file.WithPath(configFile),
	)); err != nil {
		fmt.Println(err)
		return
	}

	if err := config.Get("app").Scan(&AppConf); err != nil {
		fmt.Println(err)
		return
	}
	if AppConf.IP == "" {
		AppConf.IP = "localhost"
	}
	AppConf.Address = AppConf.IP + ":" + strconv.Itoa(AppConf.Port)

	// init logger..
	log.InitLogger(
		log.WithLevel( log.Level(AppConf.LogLevel) ),
		log.WithFields(log.Fields{
			"logger": "api",
		}),
		log.WithOutput(
			log.NewOutput(log.OutputName(AppConf.LogPath)),
		),
	)

	log.Infof("logger init, path:%v", AppConf.LogPath)

	if err := config.Get("api").Scan(&ApiConf); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ApiConf.ApiName)
}

