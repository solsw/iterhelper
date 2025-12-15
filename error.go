package iterhelper

import (
	"errors"
	"fmt"
)

var (
	ErrOddValues   = errors.New("odd number of values")
	ErrNilAction   = errors.New("nil action")
	ErrNilEqual    = errors.New("nil equal")
	ErrNilSec      = errors.New("nil Sec")
	ErrNilSec2     = errors.New("nil Sec2")
	ErrNilSelector = errors.New("nil selector")
)

func ErrWrongType(got, want any) error {
	return fmt.Errorf("wrong type: got %T, want %T", got, want)
}
