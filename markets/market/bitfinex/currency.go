package bitfinex

import (
	"github.com/cryptounicorns/trade/currencies"
)

var (
	CurrencyMapping = map[currencies.Currency]string{
		// https://api.bitfinex.com/v1/symbols
		currencies.Bitcoin:            "btc",
		currencies.Litecoin:           "ltc",
		currencies.UnitedStatesDollar: "usd",
	}
	CurrencyPairDelimiter = ""
)
