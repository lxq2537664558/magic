package config

import (
	"io/ioutil"
	"log"

	"github.com/corego/vgo/common/vlog"
	"github.com/naoina/toml"
)

type CommonConfig struct {
	Version  string
	IsDebug  bool
	LogLevel string
	LogPath  string
}

// Config ...
type Config struct {
	Common *CommonConfig
}

// Conf ...
var Conf = &Config{}

func LoadConfig() {
	initConf()
	contents, err := ioutil.ReadFile("vgo.toml")
	if err != nil {
		log.Fatal("[FATAL] load vgo.toml: ", err)
	}
	tbl, err := toml.Parse(contents)
	if err != nil {
		log.Fatal("[FATAL] parse vgo.toml: ", err)
	}
	parseCommon(tbl)
	// init log logger
	initLogger()
}

func initLogger() {
	vlog.Init(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)
}

func initConf() {
	Conf = &Config{
		Common: &CommonConfig{},
	}
}
