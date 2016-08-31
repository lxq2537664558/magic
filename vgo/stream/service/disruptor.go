package service

import disruptor "github.com/smartystreets/go-disruptor"

var controller *Controller

type Controller struct {
	controller   disruptor.Disruptor
	ring         []Metrics
	bufferMask   int64
	reservations int64
}

func NewController() *Controller {
	c := &Controller{}
	controller = c
	return c
}

func (c *Controller) Init(bufferSize int64, bufferMask int64, reservations int64) {
	c.controller = disruptor.Configure(bufferSize).
		WithConsumerGroup(Writer{}).Build()
	c.ring = make([]Metrics, bufferSize)
	c.bufferMask = bufferMask
	c.reservations = reservations
}

func (c *Controller) Start() {
	c.controller.Start()
}

func (c *Controller) Close() error {
	c.controller.Stop()
	return nil
}

func Publish(m Metrics) {
	sequence := disruptor.InitialSequenceValue
	writer := controller.controller.Writer()

	sequence = writer.Reserve(controller.reservations)
	// fmt.Println(sequence, controller.reservations)
	// 赋值
	for lower := sequence - controller.reservations + 1; lower <= sequence; lower++ {
		controller.ring[lower&controller.bufferMask] = m
	}
	// 提交
	writer.Commit(sequence-controller.reservations+1, sequence)

}
