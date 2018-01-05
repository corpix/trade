package bitfinex

import (
	"fmt"
)

type ErrSubscription struct {
	Channel    string
	Parameters interface{}
	Err        string
}

func (e *ErrSubscription) Error() string {
	return fmt.Sprintf(
		"Got an error '%s' while trying to subscribe to channel '%s' with parameters '%#v'",
		e.Err,
		e.Channel,
		e.Parameters,
	)
}

func NewErrSubscription(channel string, parameters interface{}, error string) *ErrSubscription {
	return &ErrSubscription{
		Channel:    channel,
		Parameters: parameters,
		Err:        error,
	}
}
