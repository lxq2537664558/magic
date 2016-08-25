package service



import (
	"log"
	"time"

	"github.com/corego/vgo/mecury/misc"
	"github.com/naoina/toml/ast"
	"github.com/uber-go/zap"
)

// ChainConfig chainconfig
type ChainConfig struct {
	Name   string
	Prefix string
	Suffix string

	Chain Chainer

	Interval time.Duration
}

var Chains = map[string]Chainer{}

func AddChain(name string, chain Chainer) {
	Chains[name] = chain
}

// Start init and start Chainer service
func (cc *ChainConfig) Start(stopC chan bool) {
	defer func() {
		if err := recover(); err != nil {
			misc.PrintStack(false)
			vLogger.Fatal("flush fatal error ", zap.Error(err.(error)))
		}
	}()

	cc.Chain.Init(stopC)
	go cc.Chain.Start()
}

// Show show struct message
func (cc *ChainConfig) Show() {
	log.Println("Name is ", cc.Name)
	log.Println("Prefix is ", cc.Prefix)
	log.Println("Suffix is ", cc.Suffix)
	log.Println("Interval is ", cc.Interval)
	log.Printf("Inputer is %v\n", cc.Chain)
}

type Chainer interface {
	Init(chan bool)
	Start()
	Compute(Metrics) error
}

// buildChain parses chains specific items from the ast.Table,
func buildChain(name string, tbl *ast.Table) (*ChainConfig, error) {
	ch := &ChainConfig{Name: name}

	return ch, nil
}
