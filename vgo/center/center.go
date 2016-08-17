package center

import (
	"fmt"

	"github.com/corego/vgo/vgo/config"
)

// Center struct
type Center struct {
}

// Start start center server
func (c *Center) Start() {
	fmt.Printf("config msg is %v !\n", config.Conf)

}

// Close close center server
func (c *Center) Close() error {
	return nil
}

// New get new center struct
func New() *Center {
	center := &Center{}
	return center
}
