package input

import (
	"log"

	"github.com/corego/vgo/vgo/stream"
)

// Nats nats
type Nats struct {
	Addrs  []string
	StopC  chan bool
	WriteC chan *stream.Metric
}

// NewNats return new nats
func NewNats() *Nats {
	nats := &Nats{}
	return nats
}

// Init init nats
func (n *Nats) Init(stopC chan bool, writeC chan *stream.Metric) {
	n.StopC = stopC
	n.WriteC = writeC
}

// Start start nats
func (n *Nats) Start() {
	log.Println("nats Start")
}

// Close close nats
func (n *Nats) Close() error {
	return nil
}

func init() {
	stream.AddInput("nats",
		&Nats{},
	)
}
