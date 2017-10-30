package bitfinex

import (
	"context"
	"io"
	"net/http"

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
	consumer *consumer.Consumer
	conn     io.ReadCloser
	log      loggers.Logger
}

func (c *TickerConsumer) Consume([]currencies.CurrencyPair) (<-chan *market.Ticker, error) {
	for t := range c.consumer.Consume() {
		spew.Dump(t)
	}
	panic("not impl")
	return nil, nil
}

func (c *TickerConsumer) Close() error {
	var (
		err error
	)

	err = c.consumer.Close()
	if err != nil {
		c.conn.Close()
		return err
	}

	err = c.conn.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *Bitfinex) NewTickerConsumer() (market.TickerConsumer, error) {
	var (
		l = prefixwrapper.New(
			"TickerConsumer: ",
			m.log,
		)

		r   io.ReadCloser
		c   *consumer.Consumer
		res ws.Response
		err error
	)

	r, res, err = ws.DefaultDialer.Dial(
		context.Background(),
		"wss://api.bitfinex.com/ws/2",
		http.Header{
			"Origin": []string{"http://localhost/"},
		},
	)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	c = consumer.New(
		wsutil.NewReader(
			r,
			ws.StateClientSide,
		),
		l,
	)

	return &TickerConsumer{
		consumer: c,
		conn:     r,
		log:      l,
	}, nil

}

func NewTickerConsumer() (market.TickerConsumer, error) {
	return Default.NewTickerConsumer()
}
