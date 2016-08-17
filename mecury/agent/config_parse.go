package agent

import (
	"log"
	"os"
	"time"

	"github.com/corego/vgo/mecury/misc"
	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
	"github.com/uber-go/zap"
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
			vLogger.Fatal("global_tags to table error")
		}
		err := toml.UnmarshalTable(subTbl, Conf.Tags)
		if err != nil {
			vLogger.Fatal("parse global_tags error", zap.Error(err))
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

func parseFilters(tbl *ast.Table) {
	Conf.Filter = &GlobalFilter{}
	if val, ok := tbl.Fields["global_filters"]; ok {
		if subTbl, ok := val.(*ast.Table); ok {
			if node, ok := subTbl.Fields["namedrop"]; ok {
				if kv, ok := node.(*ast.KeyValue); ok {
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								Conf.Filter.NameDrop = append(Conf.Filter.NameDrop, str.Value)
							}
						}
					}
				}
			}
		}
	}

	nameDrop, err := CompileFilter(Conf.Filter.NameDrop)
	if err != nil {
		log.Fatalf("Error compiling 'namedrop', %s\n", err)
	}

	Conf.Filter.nameDrop = nameDrop
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
