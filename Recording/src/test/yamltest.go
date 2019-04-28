package main

import (
	"config"
	"fmt"
)

func main() {
	config.InitConfig("setting.yaml")
	fmt.Println(config.AppConf.Name  == nil)
}