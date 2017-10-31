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
	for t := range c.Consumer.Consume() {
		spew.Dump(t)
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
