package bitfinex

import (
	"io"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/cryptounicorns/websocket/consumer"
	"github.com/davecgh/go-spew/spew"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/market"
)

const (
	TickerChannelName = "ticker"
)

type TickerConsumer struct {
	*consumer.Consumer

	connection io.ReadWriter
	log        loggers.Logger
}

func (c *TickerConsumer) handshake(stream <-chan []byte) error {
	var (
		e         = <-stream
		event     = &Event{}
		infoEvent = &InfoEvent{}
		err       error
	)

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

func (c *TickerConsumer) subscribe(pairs []currencies.CurrencyPair, stream <-chan []byte) error {
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
			Pair: pairs,
		},
	)
	if err != nil {
		return err
	}

	err = wsutil.WriteClientText(
		c.connection,
		e,
	)
	if err != nil {
		return err
	}

	e = <-stream

	spew.Dump(e)
	err = Format.Unmarshal(
		e,
		&event,
	)
	if err != nil {
		return err
	}

	switch event.Event {
	case SubscribedEventName:
		return nil
	case ErrorEventName:
		errorEvent := &ErrorEvent{}
		err = Format.Unmarshal(
			e,
			errorEvent,
		)
		if err != nil {
			return err
		}

		return NewErrSubscription(
			errorEvent.Channel,
			errorEvent.Msg,
		)
	default:
		return NewErrUnexpectedEvent(
			SubscribeEventName+"|"+ErrorEventName,
			event.Event,
		)
	}
}

func (c *TickerConsumer) Consume(pairs []currencies.CurrencyPair) <-chan *market.Ticker {
	func() {
		var (
			stream = c.Consumer.Consume()
			err    error
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

		err = c.subscribe(pairs, stream)
		if err != nil {
			c.log.Error(err)
			return
		}

		c.log.Print("subscribed")

		spew.Dump("we are ok!", <-stream)
		spew.Dump("we are ok2!", <-stream)
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
		connection: c,
		log:        l,
	}
}

func NewTickerConsumer(r io.ReadWriter) market.TickerConsumer {
	return Default.NewTickerConsumer(r)
}
