package market

type Market interface {
	ID() string
	GetTickers([]CurrencyPair) ([]*Ticker, error)
	GetTicker(CurrencyPair) (*Ticker, error)
	Close() error
}
