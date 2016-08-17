package stream

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
	// init the new config params
	initConf()

	contents, err := ioutil.ReadFile("vgo.toml")
	if err != nil {
		log.Fatal("[FATAL] load vgo.toml: ", err)
	}
	tbl, err := toml.Parse(contents)
	if err != nil {
		log.Fatal("[FATAL] parse vgo.toml: ", err)
	}
	// parse common config
	parseCommon(tbl)
	// init logger
	initLogger()

	// init Inputers
	parseInputs(tbl)
}

// initLogger init logger
func initLogger() {
	vlog.Init(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)
	vLogger = vlog.Logger
}

func initConf() {
	Conf = &Config{
		Common: &CommonConfig{},
	}
}
