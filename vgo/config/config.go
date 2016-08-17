package config

import (
	"io/ioutil"
	"log"

	"github.com/corego/vgo/common/vlog"

	"gopkg.in/yaml.v1"
)

// Config ...
type Config struct {
	Common struct {
		Version  string
		IsDebug  bool `yaml:"debug"`
		LogPath  string
		LogLevel string
	}
	Alarm struct {
	}
	Center struct {
	}
	Stream struct {
		Discover struct {
			Etcd struct {
				Addrs []string
			}
		}
		Inputer struct {
			Nats struct {
				Addrs []string
			}
		}
	}
}

// Conf ...
var Conf = &Config{}

func InitConf() {
	data, err := ioutil.ReadFile("vgo.yaml")
	if err != nil {
		log.Fatal("read config error :", err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Fatal("yaml decode error :", err)
	}

	// init log logger
	initLogger()
}

func initLogger() {
	vlog.Init(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)
}
