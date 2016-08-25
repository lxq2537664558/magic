package service

import "log"

// Writer send data to plugins(alarm, chain, metric)
type Writer struct {
	workn int
	recvC chan Metrics
	stopC chan bool
}

// NewWriter get new writer
func NewWriter() *Writer {
	w := &Writer{}
	return w
}

// Init init writer
func (w *Writer) Init(metricChan chan Metrics) {
	w.recvC = metricChan
	w.workn = Conf.Stream.WriterNum
	w.stopC = make(chan bool, 1)
}

// Start start write service
func (w *Writer) Start() {
	for index := 0; index < w.workn; index++ {
		num := index
		go w.Working(num)
	}
	log.Println("Writer start")
}

// Close stop writer service
func (w *Writer) Close() error {
	log.Println("Writer close")
	close(w.stopC)
	return nil
}

func (w *Writer) Working(num int) {
	// start workpool
	for {
		select {
		case data, ok := <-w.recvC:
			if ok {
				for _, c := range Conf.Alarms {
					c.Alarm.Compute(data)
				}
				for _, c := range Conf.Chains {
					c.Chain.Compute(data)
				}

				for _, c := range Conf.MetricOutputs {
					c.MetricOutput.Compute(data)
				}
				log.Println("Working number is", num, ",recv data is", data)
			}
			break
		case <-w.stopC:
			log.Println("Get stop signal!")
			return
		}
	}
}
