package markets

import (
	"strings"

	"github.com/corpix/loggers"

	"github.com/cryptounicorns/trade/markets/market"
	"github.com/cryptounicorns/trade/markets/market/bitfinex"
)

const (
	BitfinexMarket = bitfinex.Name
)

func New(market string, config Config, logger loggers.Logger) (market.Market, error) {
	switch strings.ToLower(market) {
	case BitfinexMarket:
		return bitfinex.New(
			config.Bitfinex,
			logger,
		), nil
	default:
		return nil, NewErrUnsupportedMarket(market)
	}
}
