package market

import (
	"github.com/cryptounicorns/trade/currencies"
)

type TickerConsumer interface {
	Consume([]currencies.CurrencyPair) <-chan *Ticker
	Close() error
}
