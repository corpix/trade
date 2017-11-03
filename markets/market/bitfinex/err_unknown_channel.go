package bitfinex

import (
	"fmt"
)

type ErrUnknownChannel struct {
	ChannelID uint
}

func (e *ErrUnknownChannel) Error() string {
	return fmt.Sprintf(
		"Channel with ID '%d' is not known",
		e.ChannelID,
	)
}

func NewErrUnknownChannel(channelID uint) *ErrUnknownChannel {
	return &ErrUnknownChannel{
		ChannelID: channelID,
	}
}
