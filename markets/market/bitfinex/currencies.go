package bitfinex

import (
	"github.com/cryptounicorns/trade/currencies"
)

func CurrencyPairToSymbolPair(mapper currencies.Mapper, pair currencies.CurrencyPair) (SymbolPair, error) {
	var (
		symbolPair SymbolPair
		left       currencies.Currency
		right      currencies.Currency
		err        error
	)

	left, err = mapper.ToMarket(pair.Left)
	if err != nil {
		return symbolPair, err
	}

	right, err = mapper.ToMarket(pair.Right)
	if err != nil {
		return symbolPair, err
	}

	symbolPair = left.Symbol + SymbolPairDelimiter + right.Symbol

	return symbolPair, nil
}

func CurrencyPairsToSymbolPairs(mapper currencies.Mapper, pairs []currencies.CurrencyPair) ([]SymbolPair, error) {
	var (
		symbolPairs = make([]SymbolPair, len(pairs))
		err         error
	)

	for k, v := range pairs {
		symbolPairs[k], err = CurrencyPairToSymbolPair(
			mapper,
			v,
		)
		if err != nil {
			return nil, err
		}
	}

	return symbolPairs, nil
}
