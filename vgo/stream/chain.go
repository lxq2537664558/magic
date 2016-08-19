package stream

import (
	"time"

	"github.com/naoina/toml/ast"
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

type Chainer interface {
}

// buildChain parses chains specific items from the ast.Table,
func buildChain(name string, tbl *ast.Table) (*ChainConfig, error) {
	ch := &ChainConfig{Name: name}

	return ch, nil
}
