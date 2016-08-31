package sms

import (
	"fmt"
	"time"

	"github.com/corego/vgo/vgo/stream/service"
)

type Sms struct {
	in chan *service.Alarm
}

func (c *Sms) Start() error {
	c.in = make(chan *service.Alarm, 1000)
	go func() {
		for {
			a := <-c.in

			fmt.Println("Sms Output--------------------------", time.Now())

			fmt.Println(a.Data)

			fmt.Println()
			fmt.Println()
		}
	}()
	return nil
}

func (c *Sms) Close() error {
	return nil
}

func (c *Sms) Write(a *service.Alarm) error {
	c.in <- a
	return nil
}

func init() {
	service.AddOutput("sms", &Sms{})
}
