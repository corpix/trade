package bitfinex

import (
	"fmt"
)

type ErrUnexpectedEvent struct {
	Want string
	Got  string
}

func (e *ErrUnexpectedEvent) Error() string {
	return fmt.Sprintf(
		"Unexpected event, want '%s', got '%s'",
		e.Want,
		e.Got,
	)
}

func NewErrUnexpectedEvent(w string, g string) *ErrUnexpectedEvent {
	return &ErrUnexpectedEvent{
		Want: w,
		Got:  g,
	}
}
