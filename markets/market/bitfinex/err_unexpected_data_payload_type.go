package bitfinex

import (
	"fmt"
)

type ErrUnexpectedDataPayloadType struct {
	Want string
	Got  interface{}
}

func (e *ErrUnexpectedDataPayloadType) Error() string {
	return fmt.Sprintf(
		"Unexpected data payload type, want '%s', got '%T'",
		e.Want,
		e.Got,
	)
}

func NewErrUnexpectedDataPayloadType(w string, g interface{}) *ErrUnexpectedDataPayloadType {
	return &ErrUnexpectedDataPayloadType{
		Want: w,
		Got:  g,
	}
}
