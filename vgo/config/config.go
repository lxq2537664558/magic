package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v1"
)

// Config ...
type Config struct {
	Alarm struct {
		Common struct {
			Version  string
			LogDebug bool `yaml:"debug"`
			LogPath  string
			LogLevel string
		}
	}

	Center struct {
		Common struct {
			Version  string
			LogDebug bool `yaml:"debug"`
			LogPath  string
			LogLevel string
		}
	}
	Stream struct {
		Common struct {
			Version  string
			LogDebug bool `yaml:"debug"`
			LogPath  string
			LogLevel string
		}
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

func init() {
	data, err := ioutil.ReadFile("vgo.yaml")
	if err != nil {
		log.Fatal("read config error :", err)
	}

	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		log.Fatal("yaml decode error :", err)
	}
	log.Println(Conf)
}
