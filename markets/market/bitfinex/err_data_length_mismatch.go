package bitfinex

import (
	"fmt"
)

type ErrDataLengthMismatch struct {
	Want int
	Got  int
}

func (e *ErrDataLengthMismatch) Error() string {
	return fmt.Sprintf(
		"Data length mismatch, want '%d', got '%d'",
		e.Want,
		e.Got,
	)
}

func NewErrDataLengthMismatch(w int, g int) *ErrDataLengthMismatch {
	return &ErrDataLengthMismatch{
		Want: w,
		Got:  g,
	}
}
