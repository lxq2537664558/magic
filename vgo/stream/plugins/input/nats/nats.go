package input

import (
	"fmt"

	"github.com/corego/vgo/vgo/stream"
)

// Nats nats
type Nats struct {
	Addrs []string
}

// NewNats return new nats
func NewNats() *Nats {
	nats := &Nats{}
	return nats
}

// Init init nats
func (n *Nats) Init() {

}

// Start start nats
func (n *Nats) Start() {

}

// Close close nats
func (n *Nats) Close() error {
	return nil
}

func init() {
	fmt.Println("init nats")
	stream.AddInput("nats",
		&Nats{},
	)
}
