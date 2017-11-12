package gdax

import (
	"encoding/json"
	"io"

	"github.com/corpix/loggers"
	"github.com/corpix/trade/market"
	"github.com/cryptounicorns/platypus/consumer"
)

// {
//     "type": "ticker",
//     "trade_id": 20153558,
//     "sequence": 3262786978,
//     "time": "2017-09-02T17:05:49.250000Z",
//     "product_id": "BTC-USD",
//     "price": "4388.01000000",
//     "side": "buy", // Taker side
//     "last_size": "0.03000000",
//     "best_bid": "4388",
//     "best_ask": "4388.01"
// }

type Data []json.RawMessage
type Event struct {
	Type string `json:"type"`
}

type TickerConsumer struct {
	*consumer.Consumer

	connection          io.ReadWriter
	currencies          currencies.Mapper
	channelToSymbolPair symbolPairByChannel
	tickers             chan *market.Ticker
	done                chan struct{}
	log                 loggers.Logger
}

// type Interator struct {}

func (c *TickerConsumer) subscribe(pair SymbolPair, iterator *Iterator) (uint, error) {
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
}
