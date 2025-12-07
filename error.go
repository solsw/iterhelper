package iterhelper

import (
	"errors"
)

var (
	ErrNilAction   = errors.New("nil action")
	ErrNilEqual    = errors.New("nil equal")
	ErrNilSec      = errors.New("nil Sec")
	ErrNilSec2     = errors.New("nil Sec2")
	ErrNilSelector = errors.New("nil selector")
)
