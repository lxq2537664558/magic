package center

import "fmt"

// Center struct
type Center struct {
}

// Start start center server
func (c *Center) Start() {
	fmt.Printf("Center !\n")

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
