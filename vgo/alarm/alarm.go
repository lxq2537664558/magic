package alarm

import "fmt"

// Alarm struct
type Alarm struct {
}

// Start start alarm server
func (a *Alarm) Start() {
	fmt.Printf("Alarm !\n")
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
