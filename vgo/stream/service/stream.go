package service

import (
	"fmt"
	"log"
	"time"

	"github.com/corego/vgo/vgo/stream/strategy"
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
	hosts           *strategy.Hosts
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
	s.metricChan = make(chan Metrics, 1)

	// init disruptor
	s.controller = NewController()
	s.controller.Init(Conf.Stream.DisruptorBuffersize, Conf.Stream.DisruptorBuffermask, Conf.Stream.DisruptorReservations)

	// init strategyes
	s.strategyes = strategy.NewStrategy(Conf.Stream.StrategyDbname, Conf.Stream.StrategyBucketname)
	s.strategyes.Init()

	// init alarmer
	s.alarmer = NewAlarm()
	s.alarmer.Init()

	// init hosts
	s.hosts = strategy.NewHosts()
}

func StreamTestFunc() {
	// strategy.HostTest()
	AddHost("scc@Google", "zeus")
	AddHost("scc@Google", "room")
	AddHost("scc@Google", "cache")
	AddHost("scc@Google", "center")
	AddHost("scc@Google", "vgo")
	AddHost("scc@Google", "uuid")
	gs, _ := GetGroups("scc@Google")
	go func() {
		for {
			for k, v := range gs {
				log.Println(k, v)
			}
		}
	}()

	go func() {
		for {
			for k, v := range gs {
				log.Println(k, v)
			}
		}
	}()
	time.Sleep(time.Second * 1)
	log.Println("Host get groups is ", gs)
	// DeleHost("scc@Google")
	DeleGroupInHosts("scc@Google", "zeus")
	DeleGroupInHosts("scc@Google", "room")
	DeleGroupInHosts("scc@Google", "cache")
	DeleGroupInHosts("scc@Google", "uuid")
	DeleGroupInHosts("scc@Google", "vgo")
	DeleGroupInHosts("scc@Google", "center")
	gs, _ = GetGroups("scc@Google")
	log.Println("Host get groups is ", gs)

}

func AddHost(hostname string, gid string) error {
	if streamer == nil {
		return fmt.Errorf("streamer is nil, please init stream!")
	}
	streamer.hosts.Add(hostname, gid)
	return nil
}

func GetGroups(hostname string) (map[string]bool, error) {
	if streamer == nil {
		return nil, fmt.Errorf("streamer is nil, please init stream!")
	}
	return streamer.hosts.Get(hostname), nil
}

func DeleGroupInHosts(hostname string, gid string) error {
	if streamer == nil {
		return fmt.Errorf("streamer is nil, please init stream!")
	}
	streamer.hosts.DeleGroupInHosts(hostname, gid)
	return nil
}

func DeleHost(hostname string) error {
	if streamer == nil {
		return fmt.Errorf("streamer is nil, please init stream!")
	}
	streamer.hosts.DelHost(hostname)
	return nil
}

// Start start stream server
func (s *Stream) Start(shutdown chan struct{}) {

	// StreamTestFunc()

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
