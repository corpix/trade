package markets

import (
	"github.com/cryptounicorns/trade/markets/market/bitfinex"
)

var (
	DefaultConfig = Config{
		Bitfinex: bitfinex.DefaultConfig,
	}
)

type Config struct {
	Bitfinex bitfinex.Config
}
