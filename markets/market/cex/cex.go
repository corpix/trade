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
	CurrencyMapping       = map[currencies.Currency]string{}
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

type Tickers struct {
	Data  []Ticker `json:"data"`
	Error string   `json:"error"`
}

//

func (m *Cex) ID() string { return Name }

func (m *Cex) GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	var (
		u               *url.URL
		r               *http.Response
		n               int
		pair            string
		pairs           = make(map[string]bool, len(currencyPairs))
		currencyPair    currencies.CurrencyPair
		responseTickers = &Tickers{}
		tickers         = make([]*market.Ticker, len(currencyPairs))
		ok              bool
		err             error
	)

	u, err = url.Parse(Addr)
	if err != nil {
		return nil, err
	}

	u.Path += "/tickers"

	for _, v := range currencyPairs {
		pair, err = v.Format(
			CurrencyMapping,
			CurrencyPairDelimiter,
		)
		if err != nil {
			return nil, err
		}
		pairs[v.String()] = true
		u.Path += "/" + pair
	}

	//

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
	err = json.NewDecoder(r.Body).Decode(responseTickers)
	if err != nil {
		return nil, err
	}

	if len(responseTickers.Error) > 0 {
		return nil, NewErrApi(responseTickers.Error)
	}

	n = 0
	for _, v := range responseTickers.Data {
		currencyPair, err = currencies.CurrencyPairFromString(
			v.Pair,
			CurrencyMapping,
			":",
		)
		if err != nil {
			return nil, err
		}

		if _, ok = pairs[currencyPair.String()]; !ok {
			continue
		}

		// XXX: has no avg
		tickers[n] = market.NewTicker(m, currencyPair)
		tickers[n].Low = float64(v.Low)
		tickers[n].High = float64(v.High)
		tickers[n].Last = float64(v.Last)
		tickers[n].Vol = float64(v.Vol)
		tickers[n].VolCur = float64(v.VolCur)
		tickers[n].Buy = v.Buy
		tickers[n].Sell = v.Sell
		tickers[n].Timestamp = float64(v.Timestamp)

		n++
	}

	return tickers, nil
}

func (m *Cex) GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	tickers, err := m.GetTickers(
		[]currencies.CurrencyPair{currencyPair},
	)
	if err != nil {
		return nil, err
	}

	return tickers[0], nil
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
