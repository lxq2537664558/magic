package alarm

import (
	"fmt"

	"github.com/corego/vgo/vgo/config"
)

// Alarm struct
type Alarm struct {
}

// Start start alarm server
func (a *Alarm) Start() {
	fmt.Printf("config msg is %v !\n", config.Conf.Alarm)
}

// Close close alarm server
func (a *Alarm) Close() error {
	return nil
}

// New get new alarm struct
func New() *Alarm {
	alarm := &Alarm{}
	return alarm
}
