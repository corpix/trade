package bitfinex

import (
	"github.com/cryptounicorns/trade/currencies"
)

var (
	CurrencyMapping = map[currencies.Currency]string{
		currencies.Aventus:               "AVT",
		currencies.Bitcoin:               "BTC",
		currencies.BitcoinCash:           "BCH",
		currencies.BitcoinGoldFutures:    "BTG",
		currencies.Bt1Cst:                "BT1",
		currencies.Bt2Cst:                "BT2",
		currencies.Dash:                  "DSH",
		currencies.Eidoo:                 "EDO",
		currencies.Eos:                   "EOS",
		currencies.Ethereum:              "ETH",
		currencies.EthereumClassic:       "ETC",
		currencies.Iota:                  "IOT",
		currencies.Litecoin:              "LTC",
		currencies.MetaverseEtp:          "ETP",
		currencies.Monero:                "XMR",
		currencies.Neo:                   "NEO",
		currencies.Omisego:               "OMG",
		currencies.Qtum:                  "QTM",
		currencies.Ripple:                "XRP",
		currencies.SantimentNetworkToken: "SAN",
		currencies.StreamrDatacoin:       "DAT",
		currencies.Tether:                "USDT",
		currencies.Zcash:                 "ZEC",
	}
)
