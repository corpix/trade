package bitfinex

type SubscribeTickerEvent struct {
	SubscribeEvent

	Pair SymbolPair `json:"pair"`
}
