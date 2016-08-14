package agent

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

func (o *Output) AddMetric(metric Metric) {

}

func (o *Output) Write() error {
	return nil
}
