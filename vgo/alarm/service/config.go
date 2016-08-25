package service

import (
	"io/ioutil"
	"log"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
)

var Conf *Config

type Config struct {
	Common *CommonConfig
	Nats   *NatsConfig

	Outputs map[string]*Output
}

type CommonConfig struct {
	Version  string
	IsDebug  bool
	LogLevel string
	LogPath  string
}

type NatsConfig struct {
	Addrs []string
	Topic string
}

func LoadConfig() {
	Conf = &Config{
		Common:  &CommonConfig{},
		Nats:    &NatsConfig{},
		Outputs: make(map[string]*Output),
	}

	contents, err := ioutil.ReadFile("alarm.toml")
	if err != nil {
		log.Fatal("[FATAL] load config: ", err)
	}

	tbl, err := toml.Parse(contents)
	if err != nil {
		log.Fatal("[FATAL] parse config: ", err)
	}

	parseCommon(tbl)

	parseNats(tbl)

	parseOutputs(tbl)
	for _, v := range Conf.Outputs {
		log.Println("config output ---- ", v.Name, ":", v.Output)
	}
}

func parseCommon(tbl *ast.Table) {
	if val, ok := tbl.Fields["common"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] : ", subTbl)
		}
		err := toml.UnmarshalTable(subTbl, Conf.Common)
		if err != nil {
			log.Fatalln("[FATAL] parseCommon: ", err, subTbl)
		}
	}
}

func parseNats(tbl *ast.Table) {
	if val, ok := tbl.Fields["nats"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] : ", subTbl)
		}
		err := toml.UnmarshalTable(subTbl, Conf.Nats)
		if err != nil {
			log.Fatalln("[FATAL] parseNats: ", err, subTbl)
		}
	}
}

func parseOutputs(tbl *ast.Table) {
	if val, ok := tbl.Fields["outputs"]; ok {
		subTbl, _ := val.(*ast.Table)
		for pn, pt := range subTbl.Fields {
			switch iTbl := pt.(type) {
			case *ast.Table:
				Conf.AddOutput(pn, iTbl)
			case []*ast.Table:
				for _, t := range iTbl {
					Conf.AddOutput(pn, t)
				}
			default:
				log.Fatalln("[FATAL] inputs parse error: ", iTbl)
			}
		}
	}
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

func buildOutput(name string, tbl *ast.Table) (*Output, error) {
	oc := &Output{
		Name: name,
	}

	return oc, nil
}
