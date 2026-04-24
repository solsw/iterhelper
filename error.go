package iterhelper

import (
	"errors"
)

var (
	ErrOddValues   = errors.New("odd number of values")
	ErrNilAction   = errors.New("nil action")
	ErrNilEqual    = errors.New("nil equal")
	ErrNilSeq      = errors.New("nil Seq")
	ErrNilSeq2     = errors.New("nil Seq2")
	ErrNilSelector = errors.New("nil selector")
	ErrWrongType   = errors.New("wrong type")
)
