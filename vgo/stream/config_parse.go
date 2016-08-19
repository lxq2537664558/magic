package stream

import (
	"log"

	"github.com/naoina/toml"
	"github.com/naoina/toml/ast"
	"github.com/uber-go/zap"
)

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
			if node, ok := subTbl.Fields["alarmdrop"]; ok {
				if kv, ok := node.(*ast.KeyValue); ok {
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								Conf.Filter.AlarmDrop = append(Conf.Filter.AlarmDrop, str.Value)
							}
						}
					}
				}
			}
		}
	}

	alarmDrop, err := CompileFilter(Conf.Filter.AlarmDrop)
	if err != nil {
		log.Fatalf("Error compiling 'alarmdrop', %s\n", err)
	}

	Conf.Filter.alarmDrop = alarmDrop

	// parse output plugin drop
	if val, ok := tbl.Fields["global_filters"]; ok {
		if subTbl, ok := val.(*ast.Table); ok {
			if node, ok := subTbl.Fields["metric_outputdrop"]; ok {
				if kv, ok := node.(*ast.KeyValue); ok {
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								Conf.Filter.Metric_OutputDrop = append(Conf.Filter.Metric_OutputDrop, str.Value)
							}
						}
					}
				}
			}
		}
	}

	metric_OutputDrop, err := CompileFilter(Conf.Filter.Metric_OutputDrop)
	if err != nil {
		log.Fatalf("Error compiling 'metric_outputdrop', %s\n", err)
	}

	Conf.Filter.metric_OutputDrop = metric_OutputDrop

	// parse output plugin drop
	if val, ok := tbl.Fields["global_filters"]; ok {
		if subTbl, ok := val.(*ast.Table); ok {
			if node, ok := subTbl.Fields["chaindrop"]; ok {
				if kv, ok := node.(*ast.KeyValue); ok {
					if ary, ok := kv.Value.(*ast.Array); ok {
						for _, elem := range ary.Value {
							if str, ok := elem.(*ast.String); ok {
								Conf.Filter.ChainDrop = append(Conf.Filter.ChainDrop, str.Value)
							}
						}
					}
				}
			}
		}
	}

	chainDrop, err := CompileFilter(Conf.Filter.ChainDrop)
	if err != nil {
		log.Fatalf("Error compiling 'chainDrop', %s\n", err)
	}

	Conf.Filter.chainDrop = chainDrop
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
				vLogger.Info("config", zap.String("inputer", pn))
			case []*ast.Table:
				for _, t := range iTbl {
					Conf.AddInput(pn, t)
					vLogger.Info("config", zap.String("inputer", t.Name))
				}

			default:
				log.Fatalln("[FATAL] inputs parse error: ", iTbl)
			}
		}
	}
}

func parseAlarms(tbl *ast.Table) {
	if val, ok := tbl.Fields["alarms"]; ok {
		subTbl, _ := val.(*ast.Table)
		for pn, pt := range subTbl.Fields {
			// filter the alarms,drop the ones in global_filters
			if !Conf.Filter.ShouldAlarmDropPass(pn) {
				continue
			}

			switch iTbl := pt.(type) {
			case *ast.Table:
				Conf.AddArarm(pn, iTbl)
				vLogger.Info("config", zap.String("alarmer", pn))
			case []*ast.Table:
				for _, t := range iTbl {
					Conf.AddArarm(pn, t)
					vLogger.Info("config", zap.String("alarmer", t.Name))
				}

			default:
				log.Fatalln("[FATAL] alarms parse error: ", iTbl)
			}
		}
	}
}
