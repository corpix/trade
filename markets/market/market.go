package market

import (
	"github.com/corpix/trade/currencies"
)

type Market interface {
	ID() string
	GetTickers([]currencies.CurrencyPair) ([]*Ticker, error)
	GetTicker(currencies.CurrencyPair) (*Ticker, error)
	Close() error
}
