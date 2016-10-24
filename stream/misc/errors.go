package misc

import "errors"

var (
	TimeoutErr          = errors.New("Command timed out.")
	NotImplementedError = errors.New("not implemented yet")
)
