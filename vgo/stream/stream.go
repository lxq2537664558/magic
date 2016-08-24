package stream

import (
	"log"

	"github.com/uber-go/zap"
)

var vLogger zap.Logger

type StreamConfig struct {
	InputerQueue int
	WriterNum    int
}

// Stream struct
type Stream struct {
	stopPluginsChan chan bool
	metricChan      chan *Metric
	writer          *Writer
}

// New get new stream struct
func New() *Stream {
	stream := &Stream{}
	return stream
}

// Init init stream
func (s *Stream) Init() {
	s.stopPluginsChan = make(chan bool, 1)
	s.metricChan = make(chan *Metric, Conf.Stream.InputerQueue)
	s.writer = NewWriter()
	s.writer.Init(s.metricChan)
}

// Start start stream server
func (s *Stream) Start(shutdown chan struct{}) {
	// start writer service
	s.writer.Start()

	// start plugins service
	for _, inC := range Conf.Inputs {
		inC.Start(s.stopPluginsChan, s.metricChan)
	}

	for _, amC := range Conf.Alarms {
		amC.Start(s.stopPluginsChan)
	}

	for _, chC := range Conf.Chains {
		chC.Start(s.stopPluginsChan)
	}

	for _, moC := range Conf.MetricOutputs {
		moC.Start(s.stopPluginsChan)
	}
}

// Close close stream server
func (s *Stream) Close() error {
	log.Println("Stream close!")
	close(s.stopPluginsChan)
	close(s.metricChan)

	s.writer.Close()
	return nil
}
