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

type TickerConsumer struct {
	*consumer.Consumer

	connection io.Reader
	log        loggers.Logger
}

func (c *TickerConsumer) Consume([]currencies.CurrencyPair) (<-chan *market.Ticker, error) {
	var (
		err   error
		n     uint64
		event *Event
	)
	// FIXME: After a some time it goes to the infinite loop
	// flooding the terminal with ([]uint8) (cap=512) ...
	// probably conenction is dead?
	// Yeap, I saw FIN ACK from the server after 100 seconds of incativity.
consumerLoop:
	for t := range c.Consumer.Consume() {
		n++

		// FIXME: Temporary fix for closed connection
		if len(t) == 0 {
			break
		}

		event = &Event{}
		err = Format.Unmarshal(
			t,
			event,
		)
		if err != nil {
			c.log.Errorf(
				"Got an error while unmarshal event: '%s'",
				err,
			)
			continue
		}

		switch event.Event {

		// TODO: Next:
		// - implement pings
		// - implement subscriptions

		case "info":
			// FIXME: Any chance we could make this better?
			if n != 1 {
				c.log.Error(
					NewErrUnexpectedEvent(
						event.Event,
						1,
						n,
					),
				)
				break consumerLoop
			}

			info := &InfoEvent{}

			err = Format.Unmarshal(
				t,
				info,
			)
			if err != nil {
				c.log.Error(err)
				break consumerLoop
			}

			if info.Version != Version {
				c.log.Error(
					NewErrUnsupportedAPIVersion(
						Version,
						info.Version,
					),
				)
				break consumerLoop
			}

			c.log.Print("we are ok")
			continue
		}

		spew.Dump(event)
	}
	panic("not impl")
	return nil, nil
}

func (m *Bitfinex) NewTickerConsumer(r io.Reader) market.TickerConsumer {
	var (
		l = prefixwrapper.New(
			"TickerConsumer: ",
			m.log,
		)
	)

	return &TickerConsumer{
		Consumer: consumer.New(
			wsutil.NewReader(
				r,
				ws.StateClientSide,
			),
			l,
		),
		connection: r,
		log:        l,
	}
}

func NewTickerConsumer(r io.Reader) market.TickerConsumer {
	return Default.NewTickerConsumer(r)
}
