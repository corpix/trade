package ticker

import (
	"github.com/cryptounicorns/trade/currencies"
)

// Ticker is a common ticker model which may be used across different exchanges.
type Ticker struct {
	// High is the highest price in the last 24h.
	High float64 `json:"high"`

	// Low is the lowest price in the last 24h.
	Low float64 `json:"low"`

	// Vol is a 24h volume.
	Vol float64 `json:"vol"`

	// Last price.
	Last float64 `json:"last"`

	// Bid price.
	Buy float64 `json:"buy"`

	// Ask price.
	Sell float64 `json:"sell"`

	// Timestamp is a unix-nano timestamp.
	Timestamp uint64 `json:"timestamp"`

	// SymbolPair in the common representation.
	SymbolPair currencies.SymbolPair `json:"symbolPair"`

	// Market is a cryptocurrency exchange(market) name.
	Market string `json:"market"`

	// Tags which could more precisely describe this ticker update.
	Tags []string `json:"tags"`
}
