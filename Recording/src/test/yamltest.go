package main

import (
	"config"
	"fmt"
)

func main() {
	config.InitConfig("setting.yaml")
	fmt.Println(config.GetConfig().Name  == nil)
}