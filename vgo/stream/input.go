package stream

import (
	"log"
	"time"

	"github.com/naoina/toml/ast"
)

// InputConfig inputconfig
type InputConfig struct {
	Name   string
	Prefix string
	Suffix string

	Input Inputer

	Interval time.Duration
}

func (ic *InputConfig) Show() {
	log.Println("Name is ", ic.Name)
	log.Println("Prefix is ", ic.Prefix)
	log.Println("Suffix is ", ic.Suffix)
	log.Println("Interval is ", ic.Interval)
	log.Printf("Inputer is %v\n", ic.Input)
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

// buildInput parses input specific items from the ast.Table,
func buildInput(name string, tbl *ast.Table) (*InputConfig, error) {
	cp := &InputConfig{Name: name}

	return cp, nil
}
