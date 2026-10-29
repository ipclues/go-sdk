package ipclues

import "errors"

var (
	ErrEmptyIPAddress = errors.New("ipclues: IP address must not be empty")
)
