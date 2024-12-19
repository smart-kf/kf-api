package config

import (
	"fmt"
	"github.com/make-money-fast/xconfig"
	"gopkg.in/yaml.v2"
	"sync"
)

var (
	o     sync.Once
	_conf *Config
)

func Load(file string) *Config {
	var conf Config
	err := xconfig.ParseFromFile(file, &conf)
	if err != nil {
		panic(err)
	}
	o.Do(func() {
		_conf = &conf
		data, _ := yaml.Marshal(conf)
		fmt.Println("================ config start ===========")
		fmt.Println(string(data))
		fmt.Println("================ config end ===========")
	})
	return _conf
}

func GetConfig() *Config {
	return _conf
}
