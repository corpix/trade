package bitfinex

import (
	"encoding/json"
	"strings"

	"github.com/cryptounicorns/trade/currencies"
)

type SubscribeTickerEvent struct {
	SubscribeEvent

	Pair currencies.CurrencyPair `json:"pair"`
}

type subscribeTickerEventJSON struct {
	SubscribeEvent

	Pair string `json:"pair"`
}

func (e *SubscribeTickerEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		&subscribeTickerEventJSON{
			SubscribeEvent: e.SubscribeEvent,
			Pair: currencies.CurrencyPairToString(
				e.Pair,
				CurrencyMapping,
				CurrencyPairDelimiter,
			),
		},
	)
}

func (e *SubscribeTickerEvent) UnmarshalJSON(b []byte) error {
	var (
		v    = &subscribeTickerEventJSON{}
		pair currencies.CurrencyPair
		err  error
	)

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	pair, err = currencies.CurrencyPairFromString(
		v.Pair,
		CurrencyMapping,
		CurrencyPairDelimiter,
	)
	if err != nil {
		return err
	}

	e.SubscribeEvent = v.SubscribeEvent
	e.Pair = pair

	return nil
}
