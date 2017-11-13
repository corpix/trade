package ticker

import (
	"github.com/cryptounicorns/trade/currencies"
)

type Consumer interface {
	Consume([]currencies.CurrencyPair) (<-chan Result, error)
	Close() error
}
