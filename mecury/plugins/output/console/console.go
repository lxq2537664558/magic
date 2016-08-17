package console

import (
	"fmt"
	"time"

	"github.com/corego/vgo/mecury/agent"
)

type Console struct {
}

func (c *Console) Connect() error {
	return nil
}

func (c *Console) Close() error {
	return nil
}

func (i *Console) SampleConfig() string {
	return ""
}

func (i *Console) Description() string {
	return "send metrics to console"
}

func (i *Console) Write(metrics []agent.Metric) error {
	fmt.Println("Console Output--------------------------", time.Now())
	for _, m := range metrics {
		fmt.Println(m)
	}
	fmt.Println()
	fmt.Println()
	return nil
}

func init() {
	agent.AddOutput("console", &Console{})
}
