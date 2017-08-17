package cex

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/corpix/trade/currencies"
	e "github.com/corpix/trade/errors"
	jsonTypes "github.com/corpix/trade/json"
	"github.com/corpix/trade/markets/market"
)

const (
	Name = "cex"
	Addr = "https://cex.io/api"
)

var (
	Default       market.Market
	DefaultClient = http.DefaultClient
)

var (
	CurrencyMapping       = currencies.CurrencyMapping
	CurrencyPairDelimiter = "/"
)

type Cex struct {
	client *http.Client
}

type Ticker struct {
	Pair      string                  `json:"pair"`
	High      jsonTypes.Float64String `json:"high"`
	Low       jsonTypes.Float64String `json:"low"`
	Vol       jsonTypes.Float64String `json:"volume30d"`
	VolCur    jsonTypes.Float64String `json:"volume"`
	Last      jsonTypes.Float64String `json:"last"`
	Buy       float64                 `json:"bid"`
	Sell      float64                 `json:"ask"`
	Timestamp jsonTypes.Int64String   `json:"timestamp"`
}

//

func (m *Cex) ID() string { return Name }

func (m *Cex) GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	var (
		tickers = make([]*market.Ticker, len(currencyPairs))
		err     error
	)

	for k, v := range currencyPairs {
		tickers[k], err = m.GetTicker(v)
		if err != nil {
			return nil, err
		}
	}

	return tickers, nil
}

func (m *Cex) GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	var (
		u              *url.URL
		r              *http.Response
		pair           string
		responseTicker = &Ticker{}
		ticker         *market.Ticker
		err            error
	)

	u, err = url.Parse(Addr)
	if err != nil {
		return nil, err
	}

	pair, err = currencyPair.Format(
		CurrencyMapping,
		CurrencyPairDelimiter,
	)
	if err != nil {
		return nil, err
	}

	u.Path += "/ticker/" + pair

	r, err = m.client.Get(u.String())
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, e.NewErrEndpoint(
			u.String(),
			http.StatusText(r.StatusCode),
			r.StatusCode,
			200,
		)
	}

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(responseTicker)
	if err != nil {
		return nil, err
	}

	// XXX: has no avg
	ticker = market.NewTicker(m, currencyPair)
	ticker.Low = float64(responseTicker.Low)
	ticker.High = float64(responseTicker.High)
	ticker.Last = float64(responseTicker.Last)
	ticker.Vol = float64(responseTicker.Vol)
	ticker.VolCur = float64(responseTicker.VolCur)
	ticker.Buy = responseTicker.Buy
	ticker.Sell = responseTicker.Sell
	ticker.Timestamp = float64(responseTicker.Timestamp)

	return ticker, nil
}

func (m *Cex) Close() error { return nil }

//

func GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	return Default.GetTicker(currencyPair)
}

//

func New(c *http.Client) (*Cex, error) {
	if c == nil {
		return nil, e.NewErrArgumentIsNil(c)
	}
	return &Cex{c}, nil
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
