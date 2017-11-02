package bitfinex

import (
	"encoding/json"
	"strings"

	"github.com/cryptounicorns/trade/currencies"
)

type SubscribeTickerEvent struct {
	SubscribeEvent

	Pair []currencies.CurrencyPair `json:"pair"`
}

type subscribeTickerEventJSON struct {
	SubscribeEvent

	Pair string `json:"pair"`
}

// FIXME: Map currencies to the bitfinex names!

func (e *SubscribeTickerEvent) MarshalJSON() ([]byte, error) {
	var (
		pairs = make([]string, len(e.Pair))
	)

	for k, v := range e.Pair {
		pairs[k] = v.Left.String() + v.Right.String()
	}

	return json.Marshal(
		&subscribeTickerEventJSON{
			SubscribeEvent: e.SubscribeEvent,
			Pair:           strings.Join(pairs, ","),
		},
	)
}

func (e *SubscribeTickerEvent) UnmarshalJSON(b []byte) error {
	panic("not implemented")

	var (
		v   = &subscribeTickerEventJSON{}
		err error
	)

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	//currencies.CurrencyPairFromString(s string, mapping map[currencies.Currency]string, delimiter string)
	//currencies.NewCurrencyPair(left, right)

	e.SubscribeEvent = v.SubscribeEvent
	e.Pair = nil

	return nil
}
