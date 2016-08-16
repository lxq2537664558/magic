package agent

import (
	"io/ioutil"
	"log"
	"time"

	_ "github.com/corego/vgo/common/vlog"
	"github.com/influxdata/toml"
)

type Config struct {
	Common *CommonConfig

	// default tags
	Tags map[string]string

	Agent *AgentConfig

	Inputs  []*InputConfig
	Outputs []*OutputConfig
}

var Conf *Config

func LoadConfig() {
	// init the new  config params
	initConf()

	contents, err := ioutil.ReadFile("mecury.conf")
	if err != nil {
		log.Fatal("[FATAL] load config: ", err)
	}

	tbl, err := toml.Parse(contents)
	if err != nil {
		log.Fatal("[FATAL] parse config: ", err)
	}

	// parse common config
	parseCommon(tbl)

	// parse the global tags
	parseTags(tbl)

	// parse agent
	parseAgent(tbl)

	log.Printf("%#v\n", *Conf.Common)
}

func Reload(r chan struct{}) {
	time.Sleep(5 * time.Second)
	r <- struct{}{}
}

type CommonConfig struct {
	Version  string
	IsDebug  bool   `toml:"isdebug"`
	LogLevel string `toml:"loglevel"`
	LogPath  string `toml:"logpath"`

	Hostname string
}

type AgentConfig struct {
	// Interval at which to gather information
	Interval time.Duration
	// By default, precision will be set to the same timestamp order as the
	// collection interval, with the maximum being 1s.
	//   ie, when interval = "10s", precision will be "1s"
	//       when interval = "250ms", precision will be "1ms"
	// Precision will NOT be used for service inputs. It is up to each individual
	// service input to set the timestamp at the appropriate precision.
	Precision time.Duration

	// FlushInterval is the Interval at which to flush data
	FlushInterval time.Duration

	// MetricBatchSize is the maximum number of metrics that is wrote to an
	// output plugin in one call.
	MetricBatchSize int
}

type InputConfig struct {
	Name   string
	Prefix string
	Suffix string

	Input Inputer

	Tags     map[string]string
	Filter   InputFilter
	Interval time.Duration
}

type OutputConfig struct {
	Name string

	Output Outputer

	Metrics *Buffer
}
