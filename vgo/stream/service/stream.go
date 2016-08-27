package service

import (
	"log"

	"github.com/uber-go/zap"
)

var vLogger zap.Logger

type StreamConfig struct {
	InputerQueue          int
	WriterNum             int
	DisruptorBuffersize   int64
	DisruptorBuffermask   int64
	DisruptorReservations int64
}

func (sc *StreamConfig) Show() {
	log.Println("InputerQueue", sc.InputerQueue)
	log.Println("WriterNum", sc.WriterNum)
	log.Println("DisruptorBuffersize", sc.DisruptorBuffersize)
	log.Println("DisruptorBuffermask", sc.DisruptorBuffermask)
	log.Println("DisruptorReservations", sc.DisruptorReservations)
}

// Stream struct
type Stream struct {
	stopPluginsChan chan bool
	metricChan      chan Metrics
	writer          *Writer
	controller      *Controller
}

var streamer *Stream

// New get new stream struct
func New() *Stream {
	stream := &Stream{}
	streamer = stream
	return stream
}

// Init init stream
func (s *Stream) Init() {
	s.stopPluginsChan = make(chan bool, 1)

	// init disruptor
	s.controller = NewController()
	// (bufferSize int64, bufferMask int64, reservations int64)
	s.controller.Init(Conf.Stream.DisruptorBuffersize, Conf.Stream.DisruptorBuffermask, Conf.Stream.DisruptorReservations)
}

// Start start stream server
func (s *Stream) Start(shutdown chan struct{}) {

	s.controller.Start()

	// start plugins service
	for _, c := range Conf.Inputs {
		c.Start(s.stopPluginsChan, s.metricChan)
	}

	for _, c := range Conf.Alarms {
		c.Start(s.stopPluginsChan)
	}

	for _, c := range Conf.Chains {
		c.Start(s.stopPluginsChan)
	}

	for _, c := range Conf.MetricOutputs {
		c.Start(s.stopPluginsChan)
	}
}

// Close close stream server
func (s *Stream) Close() error {
	log.Println("Stream close!")
	close(s.stopPluginsChan)
	close(s.metricChan)

	// s.writer.Close()
	s.controller.Close()
	return nil
}
