package main

import (
	"github.com/corpix/trade/market"
	"github.com/corpix/trade/markets"

	"github.com/corpix/trade/markets/bitfinex"
	//"github.com/corpix/trade/markets/btce"
	"github.com/corpix/trade/markets/dummy"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	tickers, err := markets.GetTickers(
		[]market.Market{
			dummy.Default,
			//btce.Default,
			bitfinex.Default,
		},
		[]market.CurrencyPair{
			market.NewCurrencyPair(
				market.BTC,
				market.USD,
			),
		},
	)
	spew.Dump(tickers, err)
}
