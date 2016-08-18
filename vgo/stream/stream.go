package stream

import "github.com/uber-go/zap"

var vLogger zap.Logger

// Stream struct
type Stream struct {
}

// New get new stream struct
func New() *Stream {
	LoadConfig()
	initPlugin()
	stream := &Stream{}
	return stream
}

// Init init stream
func (s *Stream) Init() {

}

// Start start stream server
func (s *Stream) Start() {
	// vLogger.Info("stream", zap.String("name", "scc"))
}

// Close close stream server
func (s *Stream) Close() error {

	return nil
}

// initPlugin loading plugin
func initPlugin() {

}
