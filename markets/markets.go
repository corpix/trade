package markets

import (
	"net/http"
	"strings"

	"github.com/cryptounicorns/trade/markets/market"
	"github.com/cryptounicorns/trade/markets/market/bitfinex"
)

const (
	BitfinexMarket = bitfinex.Name
)

func New(market string) (market.Market, error) {
	switch strings.ToLower(market) {
	case BitfinexMarket:
		// FIXME: config + logger for constructor
		return bitfinex.New()
	default:
		return nil, NewErrUnsupportedMarket(market)
	}
}
