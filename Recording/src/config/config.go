package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AppConf struct {
	Http struct {
		Address string
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
	Name string `yaml:"omitempty"`
}

var appconf *AppConf

func InitConfig(filepath string) {
	if appconf != nil {
		panic("appconf have already been init")
	}
	yamlFile, err := ioutil.ReadFile("setting.yaml")
	if err != nil {
		panic(err.Error())
	}
	appconf = &AppConf{}
	err = yaml.Unmarshal(yamlFile, &appconf)
	if err != nil {
		panic(err.Error())
	}

}

func GetConfig() *AppConf {
	if appconf == nil {
		panic("config is nil")
	}
	return appconf
}
