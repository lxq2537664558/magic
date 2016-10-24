package service

import (
	"log"

	"github.com/nats-io/nats"
)

type input struct {
	conn *nats.Conn
}

func (in *input) Start() {
	nc := initNatsConn()
	in.conn = nc

	in.conn.Subscribe(Conf.Nats.Topic, process)
}

func initNatsConn() *nats.Conn {
	opts := nats.DefaultOptions
	opts.Servers = Conf.Nats.Addrs

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
