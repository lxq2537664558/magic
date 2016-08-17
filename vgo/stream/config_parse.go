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

		log.Printf("-----------  %v\n", subTbl)
		err := toml.UnmarshalTable(subTbl, Conf.Common)
		if err != nil {
			log.Fatalln("[FATAL] parseCommon: ", err, subTbl)
		}
	}
}

func parseInputs(tbl *ast.Table) {
	if val, ok := tbl.Fields["inputs"]; ok {
		subTbl, _ := val.(*ast.Table)
		for pn, pt := range subTbl.Fields {
			// filter the inputs,drop the ones in global_filters
			// if !Conf.Filter.ShouldInputPass(pn) {
			// 	continue
			// }

			switch iTbl := pt.(type) {
			case *ast.Table:
				// Conf.AddInput(pn, iTbl)
				vLogger.Info("config", zap.String("inputer", pn))
			case []*ast.Table:
				for _, t := range iTbl {
					// Conf.AddInput(pn, t)
					vLogger.Info("config", zap.String("inputer", t.Name))
				}
			default:
				log.Fatalln("[FATAL] inputs parse error: ", iTbl)
			}
		}
	}
}
