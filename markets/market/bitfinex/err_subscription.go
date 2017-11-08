package bitfinex

import (
	"fmt"
)

type ErrSubscription struct {
	Channel string
	Err     string
}

func (e *ErrSubscription) Error() string {
	return fmt.Sprintf(
		"Got an error while trying to subscribe to channel '%s': '%s'",
		e.Channel,
		e.Err,
	)
}

func NewErrSubscription(channel string, error string) *ErrSubscription {
	return &ErrSubscription{
		Channel: channel,
		Err:     error,
	}
}
