package gdax

import (
	"io"

	"github.com/corpix/loggers"
	"github.com/cryptounicorns/platypus/consumer"
	"github.com/cryptounicorns/trade/market"
)

type TickerConsumer struct {
	*consumer.Consumer

	connection          io.ReadWriter
	currencies          currencies.Mapper
	channelToSymbolPair symbolPairByChannel
	tickers             chan *market.Ticker
	done                chan struct{}
	log                 loggers.Logger
}

func (c *TickerConsumer) preamble(pairs []SymbolPair, iterator *Iterator) error {
	var (
		handshaker = NewHandshaker(iterator, c.log)
		//	channelID uint
		err error
	)

	err = handshaker.Handshake()
	if err != nil {
		return err
	}
	c.log.Debug("Handshaked")
}

func (c *TickerConsumer) subscribe(pair SymbolPair, iterator *Iterator) (uint, error) {
	var (
		event = Event{
			Event: SubscribeEventName,
		}
		e   []byte
		err error
	)

	subscribe, err = Format.Marshal(
		&SubscribeTickerEvent{
			SubscribeEvent: SubscribeEvent{
				// TODO: hardcoded
				"ticker",
				[]string{
					"BTC-USD",
					"ETH-EUR",
				},
				[1]Channel{
					{
						Name: "ticker",
						// [] passed implicitly (yay...)
					},
				},
			},
			Pair: pair,
		},
	)
	if err != nil {
		return 0, err
	}

	err = wsutil.WriteClientText(
		c.connection,
		subscribe,
	)
	if err != nil {
		return 0, err
	}

	switch event.Type {
	case "subscribe":
		subscribedEvent := &SubscribedEvent{}
		err = Format.Unmarshal(
			e,
			subscribedEvent,
		)
		if err != nil {
			return 0, err
		}

		return subscribedEvent.ChanID, nil
	}
}
