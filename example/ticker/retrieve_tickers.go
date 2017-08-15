package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/corpix/trade/currencies"
	"github.com/corpix/trade/markets/market"
	"github.com/corpix/trade/markets/market/bitfinex"
	"github.com/corpix/trade/markets/market/dummy"
	"github.com/corpix/trade/markets/market/yobit"
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
