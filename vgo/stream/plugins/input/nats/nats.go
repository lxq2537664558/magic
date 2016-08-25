package input

import (
	"log"

	"github.com/corego/vgo/vgo/stream"
)

// Nats nats
type Nats struct {
	Addrs  []string
	StopC  chan bool
	WriteC chan stream.Metrics
}

// NewNats return new nats
func NewNats() *Nats {
	nats := &Nats{}
	return nats
}

// Init init nats
func (n *Nats) Init(stopC chan bool, writeC chan stream.Metrics) {
	n.StopC = stopC
	n.WriteC = writeC
}

// Start start nats
func (n *Nats) Start() {
	log.Println("nats Start")
}

// Close close nats
func (n *Nats) Close() error {
	close(n.StopC)
	return nil
}

func init() {
	stream.AddInput("nats",
		&Nats{},
	)
}
