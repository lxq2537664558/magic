package stream

import "time"

// InputConfig inputconfig
type InputConfig struct {
	Name   string
	Prefix string
	Suffix string

	Input Inputer

	// Tags     map[string]string
	// Filter   InputFilter
	Interval time.Duration
}

var Inputs = map[string]Inputer{}

func AddInput(name string, input Inputer) {
	Inputs[name] = input
}

type Inputer interface {
	Init()
	Start()
	Close() error
	Recv() (*Metric, error)
}
