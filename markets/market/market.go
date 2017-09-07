package market

import (
	"github.com/cryptounicorns/trade/currencies"
)

type Market interface {
	ID() string
	GetTickers([]currencies.CurrencyPair) ([]*Ticker, error)
	GetTicker(currencies.CurrencyPair) (*Ticker, error)
	Close() error
}
