package bitfinex

import (
	"github.com/cryptounicorns/trade/currencies"
)

func CurrencyPairToMarketSymbolPair(mapper currencies.Mapper, pair currencies.CurrencyPair) (SymbolPair, error) {
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

func CurrencyPairsToMarketSymbolPairs(mapper currencies.Mapper, pairs []currencies.CurrencyPair) ([]SymbolPair, error) {
	var (
		symbolPairs = make([]SymbolPair, len(pairs))
		err         error
	)

	for k, v := range pairs {
		symbolPairs[k], err = CurrencyPairToMarketSymbolPair(
			mapper,
			v,
		)
		if err != nil {
			return nil, err
		}
	}

	return symbolPairs, nil
}

func SymbolPairToCommonCurrencyPair(mapper currencies.Mapper, symbolPair SymbolPair) (currencies.CurrencyPair, error) {
	var (
		marketLeft  currencies.Currency
		marketRight currencies.Currency
		commonLeft  currencies.Currency
		commonRight currencies.Currency
		pair        currencies.CurrencyPair
		err         error
	)

	marketLeft, err = mapper.MarketBySymbol(symbolPair[0:3])
	if err != nil {
		return pair, err
	}
	marketRight, err = mapper.MarketBySymbol(symbolPair[3:])
	if err != nil {
		return pair, err
	}

	commonLeft, err = mapper.ToCommon(marketLeft)
	if err != nil {
		return pair, err
	}
	commonRight, err = mapper.ToCommon(marketRight)
	if err != nil {
		return pair, err
	}

	pair = currencies.NewCurrencyPair(
		commonLeft,
		commonRight,
	)

	return pair, nil
}

func SymbolPairToCommonSymbolPair(mapper currencies.Mapper, symbolPair SymbolPair) (currencies.SymbolPair, error) {
	var (
		marketLeft  currencies.Currency
		marketRight currencies.Currency
		commonLeft  currencies.Currency
		commonRight currencies.Currency
		pair        currencies.SymbolPair
		err         error
	)

	marketLeft, err = mapper.MarketBySymbol(symbolPair[0:3])
	if err != nil {
		return pair, err
	}
	marketRight, err = mapper.MarketBySymbol(symbolPair[3:])
	if err != nil {
		return pair, err
	}

	commonLeft, err = mapper.ToCommon(marketLeft)
	if err != nil {
		return pair, err
	}
	commonRight, err = mapper.ToCommon(marketRight)
	if err != nil {
		return pair, err
	}

	pair = currencies.NewSymbolPair(
		commonLeft.Symbol,
		commonRight.Symbol,
	)

	return pair, nil
}
