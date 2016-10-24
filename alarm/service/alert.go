package service

import (
	"sync"
	"time"
)

type Groups struct {
	groups map[string]*Group
	*sync.RWMutex
}

var gs *Groups

type Group struct {
	ID     string
	Alerts map[string]*Alert
	Users  map[string]*User
}

type Alert struct {
	Type        int32
	Value       []float64       // index 0 : total value of warn , 1 : total value of critical
	Count       []int32         //index 0: warn, 1 : critical
	NowCount    []int32         //index 0: warn, 1 : critical
	AlarmOutput []string        // warn: mail, critical: mobile
	Duration    []time.Duration // warn duration seconds, crit duration seconds
	LastTime    []time.Time     // warn: last update time, crit : last update time
}

type User struct {
	Name string
	Info map[string]string
}
