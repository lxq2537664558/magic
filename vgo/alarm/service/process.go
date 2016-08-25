package service

import "github.com/nats-io/nats"

func process(m *nats.Msg) {
	vLogger.Info(string(m.Data))
}
