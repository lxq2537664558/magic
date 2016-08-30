package service

import (
	"io/ioutil"
	"log"

	"github.com/corego/vgo/common/vlog"
	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
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
	Stream *StreamConfig

	// global filter
	Filter        *GlobalFilter
	Outputs       map[string]*Output
	Inputs        []*InputConfig
	Chains        []*ChainConfig
	MetricOutputs []*MetricOutputConfig
}

// Conf ...
var Conf = &Config{}

func LoadConfig() {
	// init the new config params
	initConf()

	contents, err := ioutil.ReadFile("stream.toml")
	if err != nil {
		log.Fatal("[FATAL] load stream.toml: ", err)
	}
	tbl, err := toml.Parse(contents)
	if err != nil {
		log.Fatal("[FATAL] parse stream.toml: ", err)
	}
	// parse common config
	parseCommon(tbl)

	// parse stream config
	parseStream(tbl)
	Conf.Stream.Show()
	// init logger
	initLogger()

	// parse global filters
	parseFilters(tbl)

	// init Inputers
	parseInputs(tbl)

	// init Outputs
	parseOutputs(tbl)

	// init Chains
	parseChains(tbl)

	// init MetricOutputs
	parseMetricOutputs(tbl)

	log.Println("All inputs ------------------------")
	for _, in := range Conf.Inputs {
		log.Println(in.Name)
	}

	log.Println("All outpus ------------------------")
	for _, out := range Conf.Outputs {
		log.Println(out.Name)
	}

	log.Println("All chains ------------------------")
	for _, out := range Conf.Chains {
		log.Println(out.Name)
	}

	log.Println("All metric_outputs ------------------------")
	for _, out := range Conf.MetricOutputs {
		log.Println(out.Name)
	}
}

// initLogger init logger
func initLogger() {
	vlog.Init(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)
	vLogger = vlog.Logger
	log.SetFlags(log.Lmicroseconds | log.Lshortfile | log.LstdFlags)
}

func initConf() {
	Conf = &Config{
		Common:  &CommonConfig{},
		Stream:  &StreamConfig{},
		Outputs: make(map[string]*Output),
		Inputs:  make([]*InputConfig, 0),
		Chains:  make([]*ChainConfig, 0),
	}
}

func (c *Config) AddInput(name string, iTbl *ast.Table) {
	input, ok := Inputs[name]
	if !ok {
		log.Fatalf("[FATAL] no plugin %v available\n", name)
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
	inC.Show()
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

	c.Outputs[name] = outC
}

func (c *Config) AddChain(name string, iTbl *ast.Table) {
	chain, ok := Chains[name]
	if !ok {
		log.Fatalf("[FATAL] no plugin %v available\n", name)
	}

	ccC, err := buildChain(name, iTbl)
	if err != nil {
		log.Fatalln("[FATAL] build chain : ", err)
	}

	err = toml.UnmarshalTable(iTbl, chain)
	if err != nil {
		log.Fatalln("[FATAL] unmarshal chain: ", err)
	}
	ccC.Chain = chain

	c.Chains = append(c.Chains, ccC)

}

func (c *Config) AddMetricOutput(name string, iTbl *ast.Table) {
	mo, ok := MetricOutputs[name]
	if !ok {
		log.Fatalf("[FATAL] no plugin %v available\n", name)
	}

	mcC, err := buildMetricOutput(name, iTbl)
	if err != nil {
		log.Fatalln("[FATAL] build MetricOutputs : ", err)
	}

	err = toml.UnmarshalTable(iTbl, mo)
	if err != nil {
		log.Fatalln("[FATAL] unmarshal MetricOutputs: ", err)
	}
	mcC.MetricOutput = mo

	c.MetricOutputs = append(c.MetricOutputs, mcC)

}
