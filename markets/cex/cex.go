package cex

import (
	"net/http"
	"net/url"

	"encoding/json"
	e "github.com/corpix/trade/errors"
	"github.com/corpix/trade/market"

	jsonTypes "github.com/corpix/trade/json"
)

const (
	Name = "cex"
	Addr = "https://cex.io/api"
)

var (
	DefaultClient = http.DefaultClient
	Default       market.Market
)

var (
	CurrencyMapping = map[market.Currency]string{
		market.BTC: "BTC",
		market.LTC: "LTC",
		market.ETH: "ETH",
		market.GHS: "GHS",
		market.USD: "USD",
		market.EUR: "EUR",
		market.RUB: "RUB",
	}
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

func (m *Cex) GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	var (
		u               *url.URL
		r               *http.Response
		n               int
		pair            string
		pairs           = make(map[string]bool, len(currencyPairs))
		currencyPair    market.CurrencyPair
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
		currencyPair, err = market.CurrencyPairFromString(
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

func (m *Cex) GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
	tickers, err := m.GetTickers(
		[]market.CurrencyPair{currencyPair},
	)
	if err != nil {
		return nil, err
	}

	return tickers[0], nil
}

func (m *Cex) Close() error { return nil }

//

func GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
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
