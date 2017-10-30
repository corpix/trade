package dummy

import (
	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/market"
)

func (m *Dummy) GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	tickers := make([]*market.Ticker, len(currencyPairs))
	for k, v := range currencyPairs {
		tickers[k] = market.NewTicker(
			m,
			v,
		)
	}

	return tickers, nil
}

func (m *Dummy) GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	return market.NewTicker(
		m,
		currencyPair,
	), nil
}

func GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	return Default.GetTicker(currencyPair)
}
