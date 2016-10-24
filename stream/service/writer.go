package service

type Writer struct{}

func (this Writer) Consume(lower, upper int64) {
	// create data pool
	m := Metrics{}
	for lower <= upper {
		m = controller.ring[lower&controller.bufferMask]
		// 消费

		streamer.alarmer.Compute(m)

		for _, c := range Conf.Chains {
			c.Chain.Compute(m)
		}

		for _, c := range Conf.MetricOutputs {
			c.MetricOutput.Compute(m)
		}

		lower++
	}
}
