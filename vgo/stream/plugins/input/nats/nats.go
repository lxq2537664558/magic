package input

import "github.com/corego/vgo/vgo/stream"

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
	// recv
	// write
}

// Recv get data from serve
func (n *Nats) Recv() (*stream.Metric, error) {
	// recv
	return nil, nil
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
