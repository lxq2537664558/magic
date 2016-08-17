package agent

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/corego/vgo/common/vlog"
	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"github.com/uber-go/zap"
)

type Config struct {
	Common *CommonConfig

	// global tags
	Tags map[string]string

	Agent *AgentConfig

	Inputs  []*InputConfig
	Outputs []*OutputConfig

	// global filter
	Filter *GlobalFilter
}

var Conf *Config
var vLogger zap.Logger

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

	// init log logger
	vlog.Init(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)
	vLogger = vlog.Logger

	// parse the global tags
	parseTags(tbl)

	// parse agent
	parseAgent(tbl)

	// parse global filters
	parseFilters(tbl)

	// parse inputs
	parseInputs(tbl)

	// parse outputs
	parseOutputs(tbl)

	vLogger.Info("config allready loaded!")
}

func Reload(r chan struct{}) {
	time.Sleep(5 * time.Second)
	r <- struct{}{}
}

type CommonConfig struct {
	Version  string
	IsDebug  bool
	LogLevel string
	LogPath  string

	Hostname string
}

func (c *Config) AddInput(name string, iTbl *ast.Table) {
	input, ok := Inputs[name]
	if !ok {
		log.Fatalf("[FATAL] no plugin %v available\n", name)
	}

	t, ok := input.(ParserInput)
	if ok {
		parser := parserInit(name, iTbl)
		t.SetParser(parser)
	}

	inC, err := buildInput(name, iTbl)
	if err != nil {
		log.Fatalln("[FATAL] build input : ", err)
	}

	err = toml.UnmarshalTable(iTbl, input)
	if err != nil {
		log.Fatalln("[FATAL] unmarshal input: ", err)
	}
	inC.Input = input

	c.Inputs = append(c.Inputs, inC)
}

func (c *Config) AddOutput(name string, iTbl *ast.Table) {
	output, ok := Outputs[name]
	if !ok {
		log.Fatalf("[FATAL] no output plugin %v available\n", name)
	}

	outC, err := buildOutput(name, iTbl)
	if err != nil {
		log.Fatalln("[FATAL] build output : ", err)
	}

	err = toml.UnmarshalTable(iTbl, output)
	if err != nil {
		log.Fatalln("[FATAL] unmarshal output: ", err)
	}
	outC.Output = output

	c.Outputs = append(c.Outputs, outC)
}
