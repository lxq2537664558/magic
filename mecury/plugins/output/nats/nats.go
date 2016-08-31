package nats

import (
	"log"
	"time"

	"github.com/corego/vgo/mecury/agent"
	"github.com/nats-io/nats"
)

type Nats struct {
	Addrs []string
	Topic string
	conn  *nats.Conn
}

//easyjson:json
type Metrics struct {
	Data     []*MetricData `json:"d"`
	Interval int           `json:"i"`
}

type MetricData struct {
	Name   string                 `json:"n"`
	Tags   map[string]string      `json:"ts"`
	Fields map[string]interface{} `json:"f"`
	Time   time.Time              `json:"t"`
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
	mData := Metrics{
		Data:     make([]*MetricData, len(metrics)),
		Interval: int(agent.Conf.Agent.Interval.Duration.Seconds()),
	}

	for i, v := range metrics {
		mData.Data[i] = &MetricData{
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
