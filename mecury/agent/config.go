package agent

import (
	"io/ioutil"
	"log"
	"time"

	_ "github.com/corego/vgo/common/vlog"
	"github.com/corego/vgo/mecury/misc"
	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
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

	// parse inputs
	parseInputs(tbl)

	// parse outputs
	parseOutputs(tbl)

	log.Printf("%#v\n", *Conf.Common)
	log.Printf("%#v\n", *Conf.Agent)
	log.Printf("%#v\n", Conf.Tags)
	for _, input := range Conf.Inputs {
		log.Printf("input %v : %#v", input.Name, input.Input)
	}

	for _, output := range Conf.Outputs {
		log.Printf("output %v : %#v", output.Name, output.Output)
	}
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

type AgentConfig struct {
	// Interval at which to gather information
	Interval misc.Duration

	// FlushInterval is the Interval at which to flush data
	FlushInterval misc.Duration

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

type OutputConfig struct {
	Name string

	Output Outputer

	Metrics *Buffer
}

// ParseConfig is a struct that covers the data types needed for all parser types,
// and can be used to instantiate _any_ of the parsers.
type ParseConfig struct {
	// Dataformat can be one of: json, influx, graphite, value, nagios
	DataFormat string

	// Separator only applied to Graphite data.
	Separator string
	// Templates only apply to Graphite data.
	Templates []string

	// TagKeys only apply to JSON data
	TagKeys []string
	// MetricName applies to JSON & value. This will be the name of the measurement.
	MetricName string

	// DataType only applies to value, this will be the type to parse value to
	DataType string

	// DefaultTags are the default tags that will be added to all parsed metrics.
	DefaultTags map[string]string
}
