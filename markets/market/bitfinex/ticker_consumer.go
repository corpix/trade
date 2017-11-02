package bitfinex

import (
	"io"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/cryptounicorns/websocket/consumer"
	//"github.com/davecgh/go-spew/spew"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/market"
)

const (
	TickerChannelName = "ticker"
)

type currencyPairByChannel map[uint]currencies.CurrencyPair

type TickerConsumer struct {
	*consumer.Consumer

	channelToCurrencyPair currencyPairByChannel
	connection            io.ReadWriter
	log                   loggers.Logger
}

// This function exists because bitfinex API is inconsistent
// as shit. It retrieves a data from the stream and checks that
// retrieved data is a hashmap, skipping arrays, which possibly could
// be received while subscribing to channels, and handles other
// shit.
func (c *TickerConsumer) nextEvent(stream <-chan []byte) ([]byte, error) {
	var (
		event []byte
	)

streamLoop:
	for {
		event = <-stream
		switch {
		case event[0] == '{':
			// Hashmap received, looks like we have a new event
			break streamLoop
		case event[0] == '[':
			// Array received, looks like we have a data
			continue streamLoop
		default:
			// Some unexpected shit is received
			// This should not happen, but WHAT IF
			return nil, NewErrUnexpectedEvent(
				"{ ... }",
				string(event),
			)
		}
	}

	return event, nil
}

func (c *TickerConsumer) handshake(stream <-chan []byte) error {
	var (
		event     = &Event{}
		infoEvent = &InfoEvent{}
		e         []byte
		err       error
	)

	e, err = c.nextEvent(stream)
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

func (c *TickerConsumer) subscribe(pair currencies.CurrencyPair, stream <-chan []byte) (uint, error) {
	var (
		event = Event{
			Event: SubscribeEventName,
		}
		e   []byte
		err error
	)

	e, err = Format.Marshal(
		&SubscribeTickerEvent{
			SubscribeEvent: SubscribeEvent{
				Event:   event,
				Channel: TickerChannelName,
			},
			Pair: pair,
		},
	)
	if err != nil {
		return 0, err
	}

	err = wsutil.WriteClientText(
		c.connection,
		e,
	)
	if err != nil {
		return 0, err
	}

	e, err = c.nextEvent(stream)
	if err != nil {
		return 0, err
	}

	err = Format.Unmarshal(
		e,
		&event,
	)
	if err != nil {
		return 0, err
	}

	switch event.Event {
	case SubscribedEventName:
		subscribedEvent := &SubscribedEvent{}
		err = Format.Unmarshal(
			e,
			subscribedEvent,
		)
		if err != nil {
			return 0, err
		}

		return subscribedEvent.ChanID, nil
	case ErrorEventName:
		errorEvent := &ErrorEvent{}
		err = Format.Unmarshal(
			e,
			errorEvent,
		)
		if err != nil {
			return 0, err
		}

		return 0, NewErrSubscription(
			errorEvent.Channel,
			errorEvent.Msg,
		)
	default:
		return 0, NewErrUnexpectedEvent(
			SubscribeEventName+"|"+ErrorEventName,
			event.Event,
		)
	}
}

func (c *TickerConsumer) Consume(pairs []currencies.CurrencyPair) <-chan *market.Ticker {
	func() {
		var (
			stream    = c.Consumer.Consume()
			channelID uint
			err       error
		)
		// FIXME: After a some time it goes to the infinite loop
		// flooding the terminal with ([]uint8) (cap=512) ...
		// probably conenction is dead?
		// Yeap, I saw FIN ACK from the server after 100 seconds of incativity.

		// FIXME: Add timeout for reading from channel

		// TODO: Next:
		// - [ ] implement pings
		// - [X] implement subscriptions
		// - [ ] transition to the "finalized" state where
		//       we only receive the ticker
		//       (and maybe dealing with pings if required)

		err = c.handshake(stream)
		if err != nil {
			c.log.Error(err)
			return
		}
		c.log.Print("handshaked")

		for _, pair := range pairs {
			channelID, err = c.subscribe(pair, stream)
			if err != nil {
				c.log.Error(err)
				return
			}

			c.channelToCurrencyPair[channelID] = pair
			c.log.Print("subscribed ", channelID, pair)
		}

		for event := range stream {
			c.log.Printf("%s", event)
		}
	}()

	panic("not going anywhere :)")
	return nil
}

// FIXME: This is shit, consumer should receive reader by semantic.
func (m *Bitfinex) NewTickerConsumer(c io.ReadWriter) market.TickerConsumer {
	var (
		l = prefixwrapper.New(
			"TickerConsumer: ",
			m.log,
		)
	)

	return &TickerConsumer{
		Consumer: consumer.New(
			wsutil.NewReader(
				c,
				ws.StateClientSide,
			),
			l,
		),
		channelToCurrencyPair: currencyPairByChannel{},
		connection:            c,
		log:                   l,
	}
}

func NewTickerConsumer(r io.ReadWriter) market.TickerConsumer {
	return Default.NewTickerConsumer(r)
}
