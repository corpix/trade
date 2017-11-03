package bitfinex

import (
	"github.com/cryptounicorns/trade/currencies"
)

type currencyPairByChannel map[uint]currencies.CurrencyPair

func (c currencyPairByChannel) Get(channelID uint) (currencies.CurrencyPair, error) {
	var (
		pair currencies.CurrencyPair
		ok   bool
	)

	pair, ok = c[channelID]
	if !ok {
		return pair, NewErrUnknownChannel(channelID)
	}

	return pair, nil
}
