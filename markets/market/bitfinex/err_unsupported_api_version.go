package bitfinex

import (
	"fmt"
)

type ErrUnsupportedAPIVersion struct {
	Want float64
	Got  float64
}

func (e *ErrUnsupportedAPIVersion) Error() string {
	return fmt.Sprintf(
		"Unsupported API version, want '%d', got '%d'",
		e.Want,
		e.Got,
	)
}

func NewErrUnsupportedAPIVersion(w float64, g float64) *ErrUnsupportedAPIVersion {
	return &ErrUnsupportedAPIVersion{
		Want: w,
		Got:  g,
	}
}
