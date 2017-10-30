package dummy

import (
	"github.com/cryptounicorns/trade/currencies"
)

var (
	CurrencyMapping = map[currencies.Currency]string{
		currencies.Bitcoin:            "btc",
		currencies.Litecoin:           "ltc",
		currencies.UnitedStatesDollar: "usd",
		currencies.Euro:               "eur",
		currencies.RussianRuble:       "rub",
	}
	CurrencyPairDelimiter = "-"
)
