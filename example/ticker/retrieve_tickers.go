package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/market"
	"github.com/cryptounicorns/trade/markets/market/bitfinex"
	"github.com/cryptounicorns/trade/markets/market/dummy"
	"github.com/cryptounicorns/trade/markets/market/yobit"
)

func main() {
	for _, v := range []market.Market{
		bitfinex.Default,
		dummy.Default,
		yobit.Default,
	} {
		tickers, err := v.GetTickers(
			[]currencies.CurrencyPair{
				currencies.NewCurrencyPair(
					currencies.Bitcoin,
					currencies.UnitedStatesDollar,
				),
				currencies.NewCurrencyPair(
					currencies.Litecoin,
					currencies.UnitedStatesDollar,
				),
			},
		)
		spew.Dump(
			v.ID(),
			tickers,
			err,
		)
	}
}
