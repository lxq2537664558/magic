package agent

import "github.com/influxdata/toml/ast"

type Outputer interface {
	// Connect to the Output
	Connect() error
	// Close any connections to the Output
	Close() error

	// Write takes in group of points to be written to the Output
	Write(metrics []Metric) error

	// Description returns a one-sentence description on the Output
	Description() string
	// SampleConfig returns the default configuration of the Output
	SampleConfig() string
}

type OutputConfig struct {
	Name string

	Output Outputer

	Metrics *Buffer
}

var Outputs = map[string]Outputer{}

func AddOutput(n string, op Outputer) {
	Outputs[n] = op
}

type Output struct {
	*OutputConfig
}

func NewOutput(name string, output Outputer) *Output {
	return &Output{
		&OutputConfig{
			Name:   name,
			Output: output,
			//Todo
			Metrics: NewBuffer(1024),
		},
	}
}

func (o *OutputConfig) AddMetric(metric Metric) {
	o.Metrics.Add(metric)
	if o.Metrics.Len() >= o.Metrics.Cap() {
		batch := o.Metrics.Batch(o.Metrics.Len())
		o.write(batch)
	}
}

func (o *OutputConfig) write(metrics []Metric) error {
	if metrics == nil || len(metrics) == 0 {
		return nil
	}

	err := o.Output.Write(metrics)
	return err
}

func (o *OutputConfig) Write() {
	batch := o.Metrics.Batch(o.Metrics.Len())
	o.write(batch)
}
func buildOutput(name string, tbl *ast.Table) (*OutputConfig, error) {
	oc := &OutputConfig{
		Name: name,
	}
	if node, ok := tbl.Fields["metric_batch_size"]; ok {
		if kv, ok := node.(*ast.KeyValue); ok {
			if str, ok := kv.Value.(*ast.Integer); ok {
				i, err := str.Int()
				if err != nil {
					return nil, err
				}
				oc.Metrics = NewBuffer(int(i))
			}
		}
	} else {
		oc.Metrics = NewBuffer(Conf.Agent.MetricBatchSize)
	}

	delete(tbl.Fields, "metric_batch_size")

	return oc, nil
}
