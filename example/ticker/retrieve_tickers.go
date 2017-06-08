package main

import (
	"github.com/corpix/trade/market"

	"github.com/corpix/trade/markets/bitfinex"
	"github.com/corpix/trade/markets/btce"
	"github.com/corpix/trade/markets/dummy"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	for _, v := range []market.Market{
		dummy.Default,
		btce.Default,
		bitfinex.Default,
	} {
		tickers, err := v.GetTickers(
			[]market.CurrencyPair{
				market.NewCurrencyPair(
					market.BTC,
					market.USD,
				),
				market.NewCurrencyPair(
					market.LTC,
					market.USD,
				),
			},
		)
		spew.Dump(tickers, err)
	}
}
