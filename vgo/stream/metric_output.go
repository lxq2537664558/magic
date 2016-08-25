package stream

import (
	"log"
	"time"

	"github.com/corego/vgo/mecury/misc"
	"github.com/naoina/toml/ast"
	"github.com/uber-go/zap"
)

// MetricOutputConfig alarmconfig
type MetricOutputConfig struct {
	Name   string
	Prefix string
	Suffix string

	MetricOutput MetricOutputer

	Interval time.Duration
}

// Start init and start MetricOutputer service
func (mc *MetricOutputConfig) Start(stopC chan bool) {
	defer func() {
		if err := recover(); err != nil {
			misc.PrintStack(false)
			vLogger.Fatal("flush fatal error ", zap.Error(err.(error)))
		}
	}()

	mc.MetricOutput.Init(stopC)
	go mc.MetricOutput.Start()
}

// Show show struct message
func (mc *MetricOutputConfig) Show() {
	log.Println("Name is ", mc.Name)
	log.Println("Prefix is ", mc.Prefix)
	log.Println("Suffix is ", mc.Suffix)
	log.Println("Interval is ", mc.Interval)
	log.Printf("Inputer is %v\n", mc.MetricOutput)
}

var MetricOutputs = map[string]MetricOutputer{}

func AddMetricOutput(name string, meto MetricOutputer) {
	MetricOutputs[name] = meto
}

type MetricOutputer interface {
	Init(chan bool)
	Start()
	Compute(Metrics) error
}

// buildMetricOutput parses MetricOutput specific items from the ast.Table,
func buildMetricOutput(name string, tbl *ast.Table) (*MetricOutputConfig, error) {
	ac := &MetricOutputConfig{Name: name}

	return ac, nil
}
