package bitfinex

import (
	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/cryptounicorns/queues/result"
)

// This Iterator thing exists because bitfinex API is inconsistent
// as shit. This Iterator retrieves a data from the stream and checks that
// retrieved data is a hashmap, skipping arrays, which possibly could
// be received while subscribing to channels, and handles other
// shit.

type Iterator struct {
	stream <-chan result.Result
	log    loggers.Logger
}

func (i *Iterator) NextEvent() ([]byte, error) {
	var (
		event result.Result
	)

streamLoop:
	for {
		event = <-i.stream

		if event.Err != nil {
			return nil, event.Err
		}

		if len(event.Value) == 0 {
			continue
		}

		switch {
		case event.Value[0] == '{':
			// Hashmap received, looks like we have a new event
			break streamLoop
		case event.Value[0] == '[':
			// Array received, looks like we have a data
			i.log.Errorf(
				"Skipping `data` while receiving `event` '%s'",
				event.Value,
			)
			continue streamLoop
		default:
			// Some unexpected shit is received
			// This should not happen, but WHAT IF
			return nil, NewErrUnexpectedEvent(
				"{ ... }",
				string(event.Value),
			)
		}
	}

	return event.Value, nil
}

func (i *Iterator) NextData() ([]byte, error) {
	var (
		data result.Result
	)

streamLoop:
	for {
		data = <-i.stream

		if data.Err != nil {
			return nil, data.Err
		}

		if len(data.Value) == 0 {
			continue
		}

		switch {
		case data.Value[0] == '[':
			// Array received, looks like we have a data
			break streamLoop
		case data.Value[0] == '{':
			// Hashmap received, looks like we have a new event
			i.log.Errorf(
				"Skipping `event` while receiving `data` '%s'",
				data.Value,
			)
			continue streamLoop
		default:
			// Some unexpected shit is received
			// This should not happen, but WHAT IF
			return nil, NewErrUnexpectedData(
				"[ ... ]",
				string(data.Value),
			)
		}
	}

	return data.Value, nil
}

func NewIterator(s <-chan result.Result, l loggers.Logger) *Iterator {
	return &Iterator{
		stream: s,
		log: prefixwrapper.New(
			"Iterator: ",
			l,
		),
	}
}
