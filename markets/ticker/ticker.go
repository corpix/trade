package ticker

import (
	"github.com/cryptounicorns/trade/currencies"
)

// Timestamp MUST be always UTC UnixNano timestamp
type Ticker struct {
	High         float64                 `json:"high"`
	Low          float64                 `json:"low"`
	Vol          float64                 `json:"vol"`
	Last         float64                 `json:"last"`
	Buy          float64                 `json:"buy"`
	Sell         float64                 `json:"sell"`
	Timestamp    uint64                  `json:"timestamp"`
	CurrencyPair currencies.CurrencyPair `json:"currencyPair"`
	Market       string                  `json:"market"`
}
