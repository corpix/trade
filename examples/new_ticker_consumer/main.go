package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/market"
	"github.com/cryptounicorns/trade/markets/market/bitfinex"
	// XXX: Import any other market
)

func main() {
	for _, v := range []market.Market{
		bitfinex.Default,
		// XXX: Append here any other market implementation
	} {
		consumer, err := v.NewTickerConsumer()
		if err != nil {
			panic(err)
		}
		defer consumer.Close()

		tickers, err := consumer.Consume(
			[]currencies.CurrencyPair{
				currencies.NewCurrencyPair(
					currencies.Bitcoin,
					currencies.UnitedStatesDollar,
				),
				// XXX: Append here any other currency pair you want to
				// get ticker for
			},
		)
		if err != nil {
			panic(err)
		}

		for ticker := range tickers {
			spew.Dump(
				v.ID(),
				ticker,
			)
		}
	}
}
