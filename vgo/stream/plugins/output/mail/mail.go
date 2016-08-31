package mail

import (
	"fmt"
	"time"

	"github.com/corego/vgo/vgo/stream/service"
)

type Mail struct {
	in chan *service.Alarm
}

func (c *Mail) Start() error {
	c.in = make(chan *service.Alarm, 1000)
	go func() {
		for {
			a := <-c.in

			fmt.Println("Mail Output--------------------------", time.Now())

			fmt.Println(a.User, ":", string(a.Data))

			fmt.Println()
			fmt.Println()
		}
	}()
	return nil
}

func (c *Mail) Close() error {
	return nil
}

func (c *Mail) Write(a *service.Alarm) error {
	c.in <- a
	return nil
}

func init() {
	service.AddOutput("mail", &Mail{})
}
