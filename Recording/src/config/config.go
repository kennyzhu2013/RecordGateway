package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strconv"
)

var AppConf struct {
	Http struct {
		Ip string
		Port int
	}
	Api struct {
		SrvName string
	}
	Ims struct {
		Ip    string
		Ports struct {
			Start int
			End   int
		}
	}
	Timeout struct {
		T1 int
		T2 int
	}
	Rabbitmq struct{
		Topic 		string
		Url   		string
		PrefetchCount int
		PrefetchGlobal bool
	}
	Name string `yaml:"omitempty"`

	// add for var
	HttpAddress string `yaml:"omitempty"`
}

func InitConfig(filepath string) {
	if "" == filepath {
		filepath = "setting.yaml"
	}
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &AppConf)
	if err != nil {
		panic(err.Error())
	}

	if "" == AppConf.Http.Ip {
		AppConf.Http.Ip = "localhost"
	}
	AppConf.HttpAddress = AppConf.Http.Ip + ":" + strconv.Itoa(AppConf.Http.Port)
}
