package stream

import (
	"fmt"

	"github.com/corego/vgo/vgo/config"
)

// Logger logger
// var Logger zap.Logger

// Stream struct
type Stream struct {
}

// Start start stream server
func (s *Stream) Start() {
	fmt.Printf("conf msg is %v \n", config.Conf)
}

// Close close stream server
func (s *Stream) Close() error {
	return nil
}

// New get new stream struct
func New() *Stream {
	stream := &Stream{}
	return stream
}
