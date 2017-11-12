package markets

import (
	"strings"

	"github.com/corpix/loggers"

	"github.com/cryptounicorns/trade/assets"
	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/market/bitfinex"
)

const (
	BitfinexMarket = bitfinex.Name
)

func New(market string, config Config, logger loggers.Logger) (Market, error) {
	var (
		loader = currencies.NewAssetLoader(assets.Asset)

		commonCurrencies currencies.Currencies
		marketCurrencies currencies.Currencies
		err              error
	)

	commonCurrencies, err = loader.Common()
	if err != nil {
		return nil, err
	}

	switch strings.ToLower(market) {
	case BitfinexMarket:
		marketCurrencies, err = loader.Market(bitfinex.Name)
		if err != nil {
			return nil, err
		}

		return bitfinex.New(
			config.Bitfinex,
			currencies.NewMapper(
				commonCurrencies,
				marketCurrencies,
			),
			logger,
		), nil
	default:
		return nil, NewErrUnsupportedMarket(market)
	}
}
