package service

import (
	"log"
	"time"

	"github.com/corego/vgo/mecury/misc"
	"github.com/naoina/toml/ast"
	"github.com/uber-go/zap"
)

// InputConfig inputconfig
type InputConfig struct {
	Name   string
	Prefix string
	Suffix string

	Input Inputer

	Interval time.Duration
}

// Start init and start Inputer service
func (ic *InputConfig) Start(stopC chan bool, writeC chan Metrics) {
	defer func() {
		if err := recover(); err != nil {
			misc.PrintStack(false)
			vLogger.Fatal("flush fatal error ", zap.Error(err.(error)))
		}
	}()

	ic.Input.Init(stopC, writeC)
	go ic.Input.Start()
}

// Show show struct message
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
	Init(chan bool, chan Metrics)
	Start()
}

// buildInput parses input specific items from the ast.Table,
func buildInput(name string, tbl *ast.Table) (*InputConfig, error) {
	cp := &InputConfig{Name: name}

	return cp, nil
}
