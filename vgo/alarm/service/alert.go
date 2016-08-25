package service

import (
	"sync"
)

type Groups struct {
	groups map[string]*Group
	*sync.RWMutex
}

var gs = &Groups{
	make(map[string]*Group),
	&sync.RWMutex{},
}

type Group struct {
	ID     string
	Alerts map[string]*Alert
	Users  []*User
}

type Alert struct {
	Type        uint8
	Value       [2]float64
	Count       []int    //index 0: warn, 1 : critical
	NowCount    []int    //index 0: warn, 1 : critical
	AlarmOutput []string // warn: mail, critical: mobile
}

type User struct {
	Info map[string]string
}
