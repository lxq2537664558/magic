package stream

import (
	"time"

	"github.com/naoina/toml/ast"
)

// AlarmConfig alarmconfig
type AlarmConfig struct {
	Name   string
	Prefix string
	Suffix string

	Alarm Alarmer

	Interval time.Duration
}

var Alarms = map[string]Alarmer{}

func AddAlarm(name string, ararm Alarmer) {
	Alarms[name] = ararm
}

type Alarmer interface {
}

// buildAlarm parses alarm specific items from the ast.Table,
func buildAlarm(name string, tbl *ast.Table) (*AlarmConfig, error) {
	ac := &AlarmConfig{Name: name}

	return ac, nil
}
