package bitfinex

import (
	"github.com/cryptounicorns/trade/currencies"
)

type Ticker [10]float64

type pairTicker struct {
	currencies.CurrencyPair
	Ticker
}
