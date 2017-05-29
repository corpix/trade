package dummy

import (
	"net/http"

	e "github.com/corpix/trade/errors"
	"github.com/corpix/trade/market"
)

const (
	Addr = "https://localhost"
)

var (
	DefaultTransport *Transport
	Default          market.Market
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

type Transport struct {
	Addr   string
	Client *http.Client
}

type Dummy struct {
	transport *Transport
}

//

func (m *Dummy) ID() string { return "dummy" }

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
	return &market.Ticker{}, nil
}

func (m *Dummy) CurrencyPair(currencyPair market.CurrencyPair) (string, error) {
	return currencyPair.Format(
		CurrencyMapping,
		CurrencyPairDelimiter,
	)
}

//

func GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
	return Default.GetTicker(currencyPair)
}

//

func NewTransport(addr string, client *http.Client) (*Transport, error) {
	if client == nil {
		return nil, e.NewErrArgumentIsNil(client)
	}

	return &Transport{
		Addr:   addr,
		Client: client,
	}, nil
}

func New(transport *Transport) (*Dummy, error) {
	if transport == nil {
		return nil, e.NewErrArgumentIsNil(transport)
	}
	return &Dummy{transport}, nil
}

//

func init() {
	var (
		err error
	)

	DefaultTransport, err = NewTransport(
		Addr,
		http.DefaultClient,
	)
	if err != nil {
		panic(err)
	}

	Default, err = New(DefaultTransport)
	if err != nil {
		panic(err)
	}
}
