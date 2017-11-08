package bitfinex

type Ticker [10]float64

type pairTicker struct {
	SymbolPair
	Ticker
}
