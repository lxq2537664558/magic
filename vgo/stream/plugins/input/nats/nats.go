package nats

import (
	"log"

	"github.com/corego/vgo/vgo/stream/service"
	"github.com/nats-io/nats"
)

// NewNats return new nats
func NewNats() *Nats {
	nats := &Nats{}
	return nats
}

var gNats *Nats

// Init init nats
func (n *Nats) Init(stopC chan bool, writeC chan service.Metrics) {
	n.StopC = stopC
	n.WriteC = writeC
}

// Start start nats
func (n *Nats) Start() {
	log.Println("nats Start")
	n.Connect()
}

// Close close nats
func (n *Nats) Close() error {
	close(n.StopC)
	return nil
}

type Nats struct {
	Addrs  []string
	Topic  string
	StopC  chan bool
	WriteC chan service.Metrics
	conn   *nats.Conn
}

func (n *Nats) Connect() error {
	nc := initNatsConn(n.Addrs, n.Topic)
	n.conn = nc
	if nc.IsClosed() == true {
		log.Fatalln("[FATAL] can't connect to nats")
	}

	gNats = n
	return nil
}

func (n *Nats) Write(data service.Metrics) {
	n.WriteC <- data
}

func dealMetricMsg(m *nats.Msg) {
	// log.Println("[DEBUG] [dealMsg] - dealMsg get msg start", m.Data)
	data := service.Metrics{}
	if err := data.UnmarshalJSON(m.Data); err != nil {
		log.Println("nats UnmarshalJSON failed, err message is ", err)
	} else {
		service.Publish(data)
		// log.Println("Nats get msg")
	}
	// gNats.Write(data)
}

func initNatsConn(addrs []string, topic string) *nats.Conn {
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

	log.Println("nats subscribe topic is", topic)
	_, err = nc.Subscribe(topic, dealMetricMsg)
	if err != nil {
		log.Fatal(" nats subscribe topic [vgo_metrics], err message is", err)
		nc.Close()
		return nil
	}
	log.Println("initNatsConn ok")

	return nc
}

func init() {
	service.AddInput("nats",
		&Nats{},
	)
}
