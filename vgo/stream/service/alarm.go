package service

import "log"

type Alarmer struct {
}

func NewAlarm() *Alarmer {
	alarmer := &Alarmer{}
	return alarmer
}

func (am *Alarmer) Init() {

}

func (am *Alarmer) Start() {

}

func (am *Alarmer) Close() error {
	return nil
}

func (am *Alarmer) Compute(m Metrics) error {
	// Compute

	// Alarm
	log.Println("Alarmer Compute message is", m)
	return nil
}
