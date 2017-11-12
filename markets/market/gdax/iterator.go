package gdax

import (
	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
)

type Iterator struct {
	stream <-chan []byte
	log    loggers.Logger
}

func (i *Iterator) NextEvent() ([]byte, error) {
	var (
		event []byte
	)

streamLoop:
	for {
		event = <-i.stream

		if len(event) == 0 {
			continue
		}

		// switch {
		// case event[0] == '{':
		//     // Hashmap received, looks like we have a new event
		//     break streamLoop
		// case event[0] == '[':
		//     // Array received, looks like we have a data
		//     i.log.Errorf(
		//         "Skipping `data` while receiving `event` '%s'",
		//         event,
		//     )
		//     continue streamLoop
		// default:
		//     // Some unexpected shit is received
		//     // This should not happen, but WHAT IF
		//     return nil, NewErrUnexpectedEvent(
		//         "{ ... }",
		//         string(event),
		//     )
		// }
	}

	return event, nil
}

func (i *Iterator) NextData() ([]byte, error) {
	var (
		data []byte
	)

streamLoop:
	for {
		data = <-i.stream

		if len(data) == 0 {
			continue
		}

		switch {
		case data[0] == '[':
			// Array received, looks like we have a data
			break streamLoop
		case data[0] == '{':
			// Hashmap received, looks like we have a new event
			i.log.Errorf(
				"Skipping `event` while receiving `data` '%s'",
				data,
			)
			continue streamLoop
		default:
			// Some unexpected shit is received
			// This should not happen, but WHAT IF
			return nil, NewErrUnexpectedData(
				"[ ... ]",
				string(data),
			)
		}
	}

	return data, nil
}

func NewIterator(s <-chan []byte, l loggers.Logger) *Iterator {
	return &Iterator{
		stream: s,
		log: prefixwrapper.New(
			"Iterator: ",
			l,
		),
	}
}
