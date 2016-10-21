package agent

import (
	"log"
	"os"
	"time"

	"github.com/aiyun/openapm/mecury/misc"
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
	// parse input plugin drop
	Conf.Filter = &GlobalFilter{}
	if val, ok := tbl.Fields["global_filters"]; ok {
		if subTbl, ok := val.(*ast.Table); ok {
			if node, ok := subTbl.Fields["inputdrop"]; ok {
				if kv, ok := node.(*ast.KeyValue); ok {
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								Conf.Filter.InputDrop = append(Conf.Filter.InputDrop, str.Value)
							}
						}
					}
				}
			}
		}
	}

	inputDrop, err := CompileFilter(Conf.Filter.InputDrop)
	if err != nil {
		log.Fatalf("Error compiling 'inputdrop', %s\n", err)
	}

	Conf.Filter.inputDrop = inputDrop

	// parse output plugin drop
	if val, ok := tbl.Fields["global_filters"]; ok {
		if subTbl, ok := val.(*ast.Table); ok {
			if node, ok := subTbl.Fields["outputdrop"]; ok {
				if kv, ok := node.(*ast.KeyValue); ok {
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								Conf.Filter.OutputDrop = append(Conf.Filter.OutputDrop, str.Value)
							}
						}
					}
				}
			}
		}
	}

	outputDrop, err := CompileFilter(Conf.Filter.OutputDrop)
	if err != nil {
		log.Fatalf("Error compiling 'outputdrop', %s\n", err)
	}

	Conf.Filter.outputDrop = outputDrop
}

func parseInputs(tbl *ast.Table) {
	if val, ok := tbl.Fields["inputs"]; ok {
		subTbl, _ := val.(*ast.Table)
		for pn, pt := range subTbl.Fields {
			// filter the inputs,drop the ones in global_filters
			if !Conf.Filter.ShouldInputPass(pn) {
				continue
			}

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
			// filter the output drop the ones in global_filters
			if !Conf.Filter.ShouldOutputPass(pn) {
				continue
			}

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
