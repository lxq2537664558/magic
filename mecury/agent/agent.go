package agent

import (
	"log"
	"sync"
	"time"

	"github.com/corego/vgo/mecury/misc"
	"github.com/uber-go/zap"
)

type Agent struct {
}

type AgentConfig struct {
	// Interval at which to gather information
	Interval misc.Duration

	// FlushInterval is the Interval at which to flush data
	FlushInterval misc.Duration

	// MetricBatchSize is the maximum number of metrics that is wrote to an
	// output plugin in one call.
	MetricBatchSize int
}

func New() *Agent {
	return &Agent{}
}

func (a *Agent) Init() {
	for _, o := range Conf.Outputs {
		err := o.Output.Connect()
		if err != nil {
			log.Fatalf("[FATAL] output %v connect error: %v \n", o.Name, err)
		}
	}
}

func (a *Agent) Start(shutdown chan struct{}) {
	var wg sync.WaitGroup

	metricC := make(chan Metric, 10000)

	// start the listener-typed inputs
	for _, input := range Conf.Inputs {
		// Start service of any ServicePlugins
		switch p := input.Input.(type) {
		case ServiceInputer:
			acc := NewAccumulate(input, metricC)

			// Service input plugins should set their own precision of their
			// metrics.
			acc.DisablePrecision()

			if err := p.Start(acc); err != nil {
				log.Fatalf("Service for input %s failed to start, exiting\n%s\n",
					input.Name, err.Error())
			}
			defer p.Stop()
		}
	}

	// round collection start time to the collection interval
	i := int64(Conf.Agent.Interval.Duration)
	time.Sleep(time.Duration(i - (time.Now().UnixNano() % i)))

	// start flusher
	wg.Add(1)
	go func() {
		a.flusher(&wg, shutdown, metricC)
	}()

	wg.Add(len(Conf.Inputs))
	for _, input := range Conf.Inputs {
		interval := Conf.Agent.Interval.Duration
		// overwrite global interval if this plugin has it's own.
		if input.Interval != 0 {
			interval = input.Interval
		}

		go func(in *InputConfig, intvl time.Duration) {
			a.gather(&wg, shutdown, in, intvl, metricC)

		}(input, interval)
	}

	wg.Wait()
}

// flusher monitors the metrics input channel and flushes on the minimum interval
func (a *Agent) flusher(wg *sync.WaitGroup, shutdown chan struct{}, metricC chan Metric) {
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			misc.PrintStack(false)
			vLogger.Fatal("flush fatal error ", zap.Error(err.(error)))
		}
	}()

	ticker := time.NewTicker(Conf.Agent.FlushInterval.Duration)
	for {
		select {
		case <-shutdown:
			a.flush()
			return
		case <-ticker.C:
			a.flush()
		case m := <-metricC:
			for _, o := range Conf.Outputs {
				o.AddMetric(m)
			}
		}
	}
}

func (a *Agent) gather(wg *sync.WaitGroup, down chan struct{}, input *InputConfig, intvl time.Duration, metricC chan Metric) {
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			misc.PrintStack(false)
			log.Printf("[FATAL] input %v error: %v\n", input.Name, err)
		}
	}()

	ticker := time.NewTicker(intvl)
	defer ticker.Stop()

	for {
		acc := NewAccumulate(input, metricC)
		acc.SetPrecision(intvl)

		err := input.Input.Gather(acc)
		if err != nil {
			vLogger.Warn("input gather error", zap.String("name", input.Name), zap.Error(err))
		}
		select {
		case <-down:
			return
		case <-ticker.C:
			continue
		}
	}
}

func (a *Agent) flush() {
	for _, o := range Conf.Outputs {
		o.Write()
	}
}
