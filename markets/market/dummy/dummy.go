package dummy

import (
	"net/http"

	"github.com/corpix/trade/currencies"
	e "github.com/corpix/trade/errors"
	"github.com/corpix/trade/markets/market"
)

const (
	Name = "dummy"
	Addr = "https://localhost"
)

var (
	Default       market.Market
	DefaultClient = http.DefaultClient
)

var (
	CurrencyMapping = map[currencies.Currency]string{
		currencies.Bitcoin:            "btc",
		currencies.Litecoin:           "ltc",
		currencies.UnitedStatesDollar: "usd",
		currencies.Euro:               "eur",
		currencies.RussianRuble:       "rub",
	}
	CurrencyPairDelimiter = "-"
)

type Dummy struct {
	client *http.Client
}

//

func (m *Dummy) ID() string { return Name }

func (m *Dummy) GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	tickers := make([]*market.Ticker, len(currencyPairs))
	for k, v := range currencyPairs {
		tickers[k] = market.NewTicker(
			m,
			v,
		)
	}
	return tickers, nil
}

func (m *Dummy) GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	return market.NewTicker(
		m,
		currencyPair,
	), nil
}

func (m *Dummy) Close() error { return nil }

//

func GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
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
