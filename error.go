package iterhelper

import (
	"errors"
)

var (
	ErrNilAction   = errors.New("nil action")
	ErrNilEqual    = errors.New("nil equal")
	ErrNilSelector = errors.New("nil selector")
	ErrNilSource   = errors.New("nil source")
)
