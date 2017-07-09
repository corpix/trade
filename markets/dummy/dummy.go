package dummy

import (
	"net/http"

	e "github.com/corpix/trade/errors"
	"github.com/corpix/trade/market"
)

const (
	Name = "dummy"
	Addr = "https://localhost"
)

var (
	DefaultClient = http.DefaultClient
	Default       market.Market
)

var (
	CurrencyMapping = map[market.Currency]string{
		market.BTC: "btc",
		market.LTC: "ltc",
		market.USD: "usd",
		market.EUR: "eur",
		market.RUB: "rub",
	}
	CurrencyPairDelimiter = "-"
)

type Dummy struct {
	client *http.Client
}

//

func (m *Dummy) ID() string { return Name }

func (m *Dummy) GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	tickers := make([]*market.Ticker, len(currencyPairs))
	for k, v := range currencyPairs {
		tickers[k] = market.NewTicker(
			m,
			v,
		)
	}
	return tickers, nil
}

func (m *Dummy) GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
	return market.NewTicker(
		m,
		currencyPair,
	), nil
}

func (m *Dummy) Close() error { return nil }

//

func GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
	return Default.GetTicker(currencyPair)
}

//

func New(c *http.Client) (*Dummy, error) {
	if c == nil {
		return nil, e.NewErrArgumentIsNil(c)
	}
	return &Dummy{c}, nil
}

//

func init() {
	var (
		err error
	)

	Default, err = New(DefaultClient)
	if err != nil {
		panic(err)
	}
}
