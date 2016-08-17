package stream

import (
	"github.com/corego/vgo/vgo/config"

	"github.com/corego/vgo/common/vlog"
	"github.com/uber-go/zap"
)

// Logger logger
// var Logger zap.Logger

// Stream struct
type Stream struct {
}

// Start start stream server
func (s *Stream) Start() {
	vlog.Logger.Info("stream", zap.String("name", "scc"))
}

// Close close stream server
func (s *Stream) Close() error {
	return nil
}

// New get new stream struct
func New() *Stream {
	initLogger()
	stream := &Stream{}
	return stream
}

// initLogger init vgo logger
func initLogger() {
	vlog.Init(config.Conf.Common.LogPath, config.Conf.Common.LogLevel, config.Conf.Common.LogDebug)
}
