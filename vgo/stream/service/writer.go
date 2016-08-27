package service

import "fmt"

type Writer struct{}

func (this Writer) Consume(lower, upper int64) {
	// create data pool
	m := Metrics{}
	for lower <= upper {
		m = controller.ring[lower&controller.bufferMask]
		// 消费
		fmt.Println("消费信息--->>> ", m)
		// ring[lower&BufferMask]
		for _, c := range Conf.Alarms {
			c.Alarm.Compute(m)
		}
		for _, c := range Conf.Chains {
			c.Chain.Compute(m)
		}

		for _, c := range Conf.MetricOutputs {
			c.MetricOutput.Compute(m)
		}

		lower++
	}
}
