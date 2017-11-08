package bitfinex

import (
	"fmt"
)

type ErrUnexpectedData struct {
	Want string
	Got  string
}

func (e *ErrUnexpectedData) Error() string {
	return fmt.Sprintf(
		"Unexpected data, want '%s', got '%s'",
		e.Want,
		e.Got,
	)
}

func NewErrUnexpectedData(w string, g string) *ErrUnexpectedData {
	return &ErrUnexpectedData{
		Want: w,
		Got:  g,
	}
}
