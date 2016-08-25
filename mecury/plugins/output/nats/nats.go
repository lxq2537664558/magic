package nats

import (
	"log"

	"github.com/corego/vgo/mecury/agent"
	"github.com/nats-io/nats"
)

type Nats struct {
	Addrs []string
	Topic string
	conn  *nats.Conn
}

func (n *Nats) Connect() error {
	nc := initNatsConn(n.Addrs)
	n.conn = nc
	if nc.IsClosed() == true {
		log.Fatalln("[FATAL] can't connect to nats")
	}

	return nil
}

func (n *Nats) Close() error {
	return nil
}

func (n *Nats) SampleConfig() string {
	return ""
}

func (n *Nats) Description() string {
	return "send metrics to nats"
}

func (n *Nats) Write(metrics []agent.Metric) error {
	mData := agent.Metrics{
		Data: make([]*agent.MetricData, len(metrics)),
	}

	for i, v := range metrics {
		mData.Data[i] = &agent.MetricData{
			Name:   v.Name(),
			Tags:   v.Tags(),
			Fields: v.Fields(),
			Time:   v.Time(),
		}
	}

	b, err := mData.MarshalJSON()
	if err != nil {
		log.Println("[WARN] data to nats Marshal error :", err)
	}

	err = n.conn.Publish(n.Topic, b)
	if err != nil {
		log.Println("[WARN] nats publish error: ", err)
	}

	return nil
}

func init() {
	agent.AddOutput("nats", &Nats{})
}

func initNatsConn(addrs []string) *nats.Conn {
	opts := nats.DefaultOptions
	opts.Servers = addrs

	nc, err := opts.Connect()
	if err != nil {
		log.Fatal("[FATAL] init nats producer error: ", err)
	}

	// Setup callbacks to be notified on disconnects and reconnects
	nc.Opts.DisconnectedCB = func(nc *nats.Conn) {
		log.Printf("%v got disconnected!\n", nc.ConnectedUrl())
	}

	// See who we are connected to on reconnect.
	nc.Opts.ReconnectedCB = func(nc *nats.Conn) {
		log.Printf("Got reconnected to %v!\n", nc.ConnectedUrl())
	}

	return nc
}
