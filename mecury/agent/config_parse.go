package agent

import (
	"log"
	"time"

	"github.com/influxdata/toml"
	"github.com/influxdata/toml/ast"
)

func initConf() {
	Conf = &Config{
		// Agent defaults:
		Agent: &AgentConfig{
			Interval:        10 * time.Second,
			FlushInterval:   10 * time.Second,
			Precision:       10 * time.Second,
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
	if val, ok := tbl.Fields["tags"]; ok {
		subTbl, ok := val.(*ast.Table)
		if !ok {
			log.Fatalln("[FATAL] : ", subTbl)
		}
		err := toml.UnmarshalTable(subTbl, Conf.Tags)
		if err != nil {
			log.Fatalln("[FATAL] parseTags: ", err)
		}
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
