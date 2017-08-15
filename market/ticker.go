package market

type Ticker struct {
	High         float64      `json:"high"`
	Low          float64      `json:"low"`
	Avg          float64      `json:"avg"`
	Vol          float64      `json:"vol"`
	VolCur       float64      `json:"volCur"`
	Last         float64      `json:"last"`
	Buy          float64      `json:"buy"`
	Sell         float64      `json:"sell"`
	Timestamp    float64      `json:"timestamp"`
	CurrencyPair CurrencyPair `json:"currencyPair"`
	Market       string       `json:"market"`
}

func NewTicker(market Market, pair CurrencyPair) *Ticker {
	return &Ticker{
		Market:       market.ID(),
		CurrencyPair: pair,
	}
}
