package service

import "github.com/naoina/toml/ast"

type Outputer interface {
	// Connect to the Output
	Start() error
	// Close any connections to the Output
	Close() error

	// Write takes in group of points to be written to the Output
	Write(*Alarm) error
}

type Output struct {
	Name string

	Output Outputer
}

type Alarm struct {
	Data []byte
	User string
}

func (o *Output) Write(alarm *Alarm) {
	o.Output.Write(alarm)
}

var Outputs = map[string]Outputer{}

func AddOutput(n string, op Outputer) {
	Outputs[n] = op
}

func buildOutput(name string, tbl *ast.Table) (*Output, error) {
	oc := &Output{
		Name: name,
	}

	return oc, nil
}
