package btce

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/corpix/trade/currencies"
	e "github.com/corpix/trade/errors"
	"github.com/corpix/trade/markets/market"
)

const (
	Name = "btce"
	Addr = "https://btc-e.com/api/3"
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
		currencies.RussianRuble:       "rur",
	}
	CurrencyPairDelimiter = "_"
)

type Btce struct {
	client *http.Client
}

type Ticker struct {
	High    float64 `json:"high"`
	Low     float64 `json:"low"`
	Avg     float64 `json:"avg"`
	Vol     float64 `json:"vol"`
	VolCur  float64 `json:"vol_cur"`
	Last    float64 `json:"last"`
	Buy     float64 `json:"buy"`
	Sell    float64 `json:"sell"`
	Updated int64   `json:"updated"`
}

//

func (m *Btce) ID() string { return Name }

func (m *Btce) GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	var (
		u               *url.URL
		r               *http.Response
		n               int
		pair            currencies.CurrencyPair
		pairs           = make([]string, len(currencyPairs))
		responseTickers = make(map[string]Ticker, len(currencyPairs))
		tickers         = make([]*market.Ticker, len(currencyPairs))
		err             error
	)

	u, err = url.Parse(Addr)
	if err != nil {
		return nil, err
	}

	u.Path += "/ticker"

	for k, v := range currencyPairs {
		pairs[k], err = v.Format(
			CurrencyMapping,
			CurrencyPairDelimiter,
		)
		if err != nil {
			return nil, err
		}
	}
	u.Path += "/" + strings.Join(
		pairs,
		"-",
	)

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
	err = json.NewDecoder(r.Body).Decode(&responseTickers)
	if err != nil {
		return nil, err
	}

	n = 0
	for k, v := range responseTickers {
		pair, err = currencies.CurrencyPairFromString(
			k,
			CurrencyMapping,
			CurrencyPairDelimiter,
		)
		if err != nil {
			return nil, err
		}

		tickers[n] = market.NewTicker(m, pair)
		tickers[n].Avg = v.Avg
		tickers[n].Buy = v.Buy
		tickers[n].High = v.High
		tickers[n].Last = v.Last
		tickers[n].Low = v.Low
		tickers[n].Sell = v.Sell
		tickers[n].Timestamp = float64(v.Updated)
		tickers[n].Vol = v.Vol
		tickers[n].VolCur = v.VolCur
		n++
	}

	return tickers, nil
}

func (m *Btce) GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	tickers, err := m.GetTickers(
		[]currencies.CurrencyPair{currencyPair},
	)
	if err != nil {
		return nil, err
	}

	return tickers[0], nil
}

func (m *Btce) Close() error { return nil }

//

func GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	return Default.GetTicker(currencyPair)
}

//

func New(c *http.Client) (*Btce, error) {
	if c == nil {
		return nil, e.NewErrArgumentIsNil(c)
	}
	return &Btce{c}, nil
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
