package main

import (
	"io"

	"github.com/davecgh/go-spew/spew"

	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/market"
	"github.com/cryptounicorns/trade/markets/market/bitfinex"
	// XXX: Import any other market
)

func main() {
	var (
		connection io.ReadWriteCloser
		consumer   market.TickerConsumer
		tickers    <-chan *market.Ticker
		err        error
	)

	for _, v := range []market.Market{
		bitfinex.Default,
		// XXX: Append here any other market implementation
	} {
		connection, err = v.Connect()
		if err != nil {
			panic(err)
		}
		defer connection.Close()

		consumer = v.NewTickerConsumer(connection)
		defer consumer.Close()

		tickers, err = consumer.Consume(
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
