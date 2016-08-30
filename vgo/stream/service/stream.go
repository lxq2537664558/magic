package service

import (
	"alert/strategy"
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
	StrategyDbname        string
	StrategyBucketname    string
}

func (sc *StreamConfig) Show() {
	log.Println("InputerQueue", sc.InputerQueue)
	log.Println("WriterNum", sc.WriterNum)
	log.Println("DisruptorBuffersize", sc.DisruptorBuffersize)
	log.Println("DisruptorBuffermask", sc.DisruptorBuffermask)
	log.Println("DisruptorReservations", sc.DisruptorReservations)
	log.Println("StrategyDbName", sc.StrategyDbname)
	log.Println("StrategyBucketName", sc.StrategyBucketname)
}

// Stream struct
type Stream struct {
	stopPluginsChan chan bool
	metricChan      chan Metrics
	writer          *Writer
	controller      *Controller
	strategyes      *strategy.Strategy
	alarmer         *Alarmer
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
	s.controller.Init(Conf.Stream.DisruptorBuffersize, Conf.Stream.DisruptorBuffermask, Conf.Stream.DisruptorReservations)

	// init strategyes
	s.strategyes = strategy.NewStrategy(Conf.Stream.StrategyDbname, Conf.Stream.StrategyBucketname)
	s.strategyes.Init()

	// init alarmer
	s.alarmer = NewAlarm()
	s.alarmer.Init()
}

// Start start stream server
func (s *Stream) Start(shutdown chan struct{}) {

	s.controller.Start()

	s.alarmer.Start()

	// start plugins service
	for _, c := range Conf.Inputs {
		c.Start(s.stopPluginsChan, s.metricChan)
	}

	for _, c := range Conf.Outputs {
		if err := c.Output.Start(); err != nil {
			log.Fatal("Output ", c.Name, " Start failed, err message is", err)
		}
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
	s.alarmer.Close()
	return nil
}
