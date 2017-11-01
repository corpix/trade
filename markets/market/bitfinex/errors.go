package bitfinex

import (
	"fmt"
)

type ErrUnexpectedEvent struct {
	Event string
	Want  uint64
	Got   uint64
}

func (e *ErrUnexpectedEvent) Error() string {
	return fmt.Sprintf(
		"Unexpected event '%s', wanted this event to be '%d' in sequence, but it is '%d'",
		e.Event,
		e.Want,
		e.Got,
	)
}

func NewErrUnexpectedEvent(e string, w uint64, g uint64) *ErrUnexpectedEvent {
	return &ErrUnexpectedEvent{
		Event: e,
		Want:  w,
		Got:   g,
	}
}

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
