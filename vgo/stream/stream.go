package stream

import (
	"log"

	"github.com/uber-go/zap"
)

var vLogger zap.Logger

// Stream struct
type Stream struct {
	stopPluginsChan chan bool
	metricChan      chan *Metric
}

// New get new stream struct
func New() *Stream {
	stream := &Stream{
		stopPluginsChan: make(chan bool, 1),
		metricChan:      make(chan *Metric, Conf.Common.InputerQueue),
	}
	return stream
}

// Init init stream
func (s *Stream) Init() {

}

// Start start stream server
func (s *Stream) Start(shutdown chan struct{}) {
	for _, inC := range Conf.Inputs {
		inC.Start(s.stopPluginsChan, s.metricChan)
	}

	for _, amC := range Conf.Alarms {
		amC.Start(s.stopPluginsChan)
	}
}

// Close close stream server
func (s *Stream) Close() error {
	log.Println("Stream close!")
	close(s.stopPluginsChan)
	close(s.metricChan)

	return nil
}
