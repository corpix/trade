package main

import (
	"io"

	"github.com/corpix/loggers/logger/logrus"
	"github.com/davecgh/go-spew/spew"
	logrusLogger "github.com/sirupsen/logrus"

	"github.com/cryptounicorns/trade/assets"
	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/market"
	"github.com/cryptounicorns/trade/markets/market/gdax"
	// XXX: Import any other market
)

func configuredLogrus() *logrusLogger.Logger {
	var (
		l = logrusLogger.New()
	)

	l.Level = logrusLogger.DebugLevel

	return l
}

func main() {
	var (
		connection io.ReadWriteCloser
		consumer   market.TickerConsumer
		tickers    <-chan *market.Ticker

		loader           = currencies.NewAssetLoader(assets.Asset)
		commonCurrencies currencies.Currencies
		marketCurrencies currencies.Currencies
		mapper           currencies.Mapper
		market           market.Market

		// bitcoin currencies.Currency
		// dollar  currencies.Currency

		log = logrus.New(configuredLogrus())

		err error
	)

	// commonCurrencies, err = loader.Common()
	// if err != nil {
	//     panic(err)
	// }
	// marketCurrencies, err = loader.Market(bitfinex.Name)
	// if err != nil {
	//     panic(err)
	// }

	// mapper = currencies.NewMapper(
	//     commonCurrencies,
	//     marketCurrencies,
	// )

	market = gdax.New(
		gdax.DefaultConfig,
		mapper,
		log,
	)

	// bitcoin, err = mapper.CommonByName("bitcoin")
	// if err != nil {
	//     panic(err)
	// }

	// dollar, err = mapper.CommonByName("united-states-dollar")
	// if err != nil {
	//     panic(err)
	// }

	connection, err = market.Connect()
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	consumer = market.NewTickerConsumer(connection)
	defer consumer.Close()

	tickers = consumer.Consume(
		[]currencies.CurrencyPair{
			currencies.NewCurrencyPair(
				bitcoin,
				dollar,
			),
		},
	)
	if err != nil {
		panic(err)
	}

	for ticker := range tickers {
		spew.Dump(
			market.Name(),
			ticker,
		)
	}
}
