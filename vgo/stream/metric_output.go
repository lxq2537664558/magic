package stream

import (
	"time"

	"github.com/naoina/toml/ast"
)

// MetricOutputConfig alarmconfig
type MetricOutputConfig struct {
	Name   string
	Prefix string
	Suffix string

	MetricOutput MetricOutputer

	Interval time.Duration
}

var MetricOutputs = map[string]MetricOutputer{}

func AddMetricOutput(name string, meto MetricOutputer) {
	MetricOutputs[name] = meto
}

type MetricOutputer interface {
}

// buildMetricOutput parses MetricOutput specific items from the ast.Table,
func buildMetricOutput(name string, tbl *ast.Table) (*MetricOutputConfig, error) {
	ac := &MetricOutputConfig{Name: name}

	return ac, nil
}
