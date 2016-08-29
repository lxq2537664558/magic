package service

import (
	"fmt"

	"github.com/corego/vgo/common/vlog"
	"github.com/uber-go/zap"
)

var vLogger zap.Logger

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (a *Service) Start() {
	LoadConfig()
	// init log logger
	vlog.Init(Conf.Common.LogPath, Conf.Common.LogLevel, Conf.Common.IsDebug)
	vLogger = vlog.Logger

	vLogger.Info(fmt.Sprintf("config: %v", Conf))

	// init input
	input := &input{}
	input.Start()

	// init output
	for _, o := range Conf.Outputs {
		o.Output.Start()
	}

	startManager()
}
