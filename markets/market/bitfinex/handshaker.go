package bitfinex

import (
	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
)

type Handshaker struct {
	iterator *Iterator
	log      loggers.Logger
}

func (h *Handshaker) Handshake() error {
	var (
		event     = &Event{}
		infoEvent = &InfoEvent{}
		e         []byte
		err       error
	)

	e, err = h.iterator.NextEvent()
	if err != nil {
		return err
	}

	err = Format.Unmarshal(
		e,
		event,
	)
	if err != nil {
		return err
	}

	if event.Event != InfoEventName {
		return NewErrUnexpectedEvent(
			InfoEventName,
			event.Event,
		)
	}

	err = Format.Unmarshal(
		e,
		infoEvent,
	)
	if err != nil {
		return err
	}

	if infoEvent.Version != Version {
		return NewErrUnsupportedAPIVersion(
			Version,
			infoEvent.Version,
		)
	}

	return nil
}

func NewHandshaker(i *Iterator, l loggers.Logger) *Handshaker {
	return &Handshaker{
		iterator: i,
		log: prefixwrapper.New(
			"Handshaker: ",
			l,
		),
	}
}
