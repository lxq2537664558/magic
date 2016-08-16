package agent

import (
	"log"
	"sync"
	"time"

	"github.com/sunface/tools"
)

type Agent struct {
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

	// round collection start time to the collection interval
	i := int64(Conf.Agent.Interval.Duration)
	time.Sleep(time.Duration(i - (time.Now().UnixNano() % i)))

	// start flusher
	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
			if err := recover(); err != nil {
				tools.PrintStack(false)
				log.Println("[FATAL] flusher error: ", err)
			}
		}()

		a.flush(shutdown, metricC)
	}()

	wg.Add(len(Conf.Inputs))
	for _, input := range Conf.Inputs {
		interval := Conf.Agent.Interval.Duration
		// overwrite global interval if this plugin has it's own.
		if input.Interval != 0 {
			interval = input.Interval
		}

		go func(in *InputConfig, intvl time.Duration) {
			defer func() {
				wg.Done()
				if err := recover(); err != nil {
					tools.PrintStack(false)
					log.Printf("[FATAL] input %v error: %v\n", input.Name, err)
				}
			}()

			a.gather(shutdown, in, intvl, metricC)

		}(input, interval)
	}

	wg.Wait()
}

// flusher monitors the metrics input channel and flushes on the minimum interval
func (a *Agent) flush(shutdown chan struct{}, metricC chan Metric) {

}

func (a *Agent) gather(down chan struct{}, input *InputConfig, intvl time.Duration, metricC chan Metric) error {

	return nil
}
