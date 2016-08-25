package nats

import (
	"log"

	"github.com/corego/vgo/vgo/stream/service"
)

// Nats nats
type Nats struct {
	Addrs  []string
	StopC  chan bool
	WriteC chan service.Metrics
}

// NewNats return new nats
func NewNats() *Nats {
	nats := &Nats{}
	return nats
}

// Init init nats
func (n *Nats) Init(stopC chan bool, writeC chan service.Metrics) {
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
	service.AddInput("nats",
		&Nats{},
	)
}
