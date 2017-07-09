package market

type Ticker struct {
	High      float64
	Low       float64
	Avg       float64
	Vol       float64
	VolCur    float64
	Last      float64
	Buy       float64
	Sell      float64
	Timestamp float64

	CurrencyPair CurrencyPair
	Market       string
}

func NewTicker(market Market, pair CurrencyPair) *Ticker {
	return &Ticker{
		Market:       market.ID(),
		CurrencyPair: pair,
	}
}
