package service

import (
	"fmt"

	"github.com/corego/vgo/common/vlog"
	"github.com/corego/vgo/vgo/alarm/config"
	"github.com/uber-go/zap"
)

var vLogger zap.Logger

type Alarm struct {
}

func NewAlarm() *Alarm {
	return &Alarm{}
}

func (a *Alarm) Start() {
	// init log logger
	vlog.Init(config.Conf.Common.LogPath, config.Conf.Common.LogLevel, config.Conf.Common.IsDebug)
	vLogger = vlog.Logger

	vLogger.Info(fmt.Sprintf("config: %v", config.Conf))

	// init input
	input := &input{}
	input.Start()
}
