package stream

import (
	"log"
	"time"

	"github.com/corego/vgo/mecury/misc"
	"github.com/naoina/toml/ast"
	"github.com/uber-go/zap"
)

// AlarmConfig alarmconfig
type AlarmConfig struct {
	Name   string
	Prefix string
	Suffix string

	Alarm Alarmer

	Interval time.Duration
}

// Start init and start Alarmer service
func (amc *AlarmConfig) Start(stopC chan bool) {
	defer func() {
		if err := recover(); err != nil {
			misc.PrintStack(false)
			vLogger.Fatal("flush fatal error ", zap.Error(err.(error)))
		}
	}()

	amc.Alarm.Init(stopC)
	go amc.Alarm.Start()
}

// Show show struct message
func (amc *AlarmConfig) Show() {
	log.Println("Name is ", amc.Name)
	log.Println("Prefix is ", amc.Prefix)
	log.Println("Suffix is ", amc.Suffix)
	log.Println("Interval is ", amc.Interval)
	log.Printf("Inputer is %v\n", amc.Alarm)
}

var Alarms = map[string]Alarmer{}

func AddAlarm(name string, ararm Alarmer) {
	Alarms[name] = ararm
}

type Alarmer interface {
	Init(chan bool)
	Start()
	Compute(*Metric) error
}

// buildAlarm parses alarm specific items from the ast.Table,
func buildAlarm(name string, tbl *ast.Table) (*AlarmConfig, error) {
	ac := &AlarmConfig{Name: name}

	return ac, nil
}
