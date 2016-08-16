package agent

import (
	"log"
	"os"
	"time"

	"github.com/corego/vgo/mecury/misc"
	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
)

func initConf() {
	Conf = &Config{
		// Agent defaults:
		Agent: &AgentConfig{
			Interval:        misc.Duration{10 * time.Second},
			FlushInterval:   misc.Duration{10 * time.Second},
			MetricBatchSize: 1000,
		},

		Tags:    make(map[string]string),
		Inputs:  make([]*InputConfig, 0),
		Outputs: make([]*OutputConfig, 0),
		Common:  &CommonConfig{},
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

func parseTags(tbl *ast.Table) {
	if val, ok := tbl.Fields["global_tags"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] : ", subTbl)
		}
		err := toml.UnmarshalTable(subTbl, Conf.Tags)
		if err != nil {
			log.Fatalln("[FATAL] parseTags: ", err)
		}
	}

	// 解析hostname
	var host string
	var err error
	if Conf.Common.Hostname == "" {
		host, err = os.Hostname()
		if err != nil {
			log.Fatalln("[FATAL] get hostname error: ", err)
		}
		Conf.Tags["host"] = host
	} else {
		Conf.Tags["host"] = Conf.Common.Hostname
	}

}

func parseAgent(tbl *ast.Table) {
	if val, ok := tbl.Fields["agent"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] : ", subTbl)
		}
		err := toml.UnmarshalTable(subTbl, Conf.Agent)
		if err != nil {
			log.Fatalln("[FATAL] parseAgent: ", err)
		}
	}
}

func parseInputs(tbl *ast.Table) {
	if val, ok := tbl.Fields["inputs"]; ok {
		subTbl, _ := val.(*ast.Table)
		for pn, pt := range subTbl.Fields {
			switch iTbl := pt.(type) {
			case *ast.Table:
				Conf.AddInput(pn, iTbl)
			case []*ast.Table:
				for _, t := range iTbl {
					Conf.AddInput(pn, t)
				}
			default:
				log.Fatalln("[FATAL] inputs parse error: ", iTbl)
			}
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
