package bitfinex

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/cryptounicorns/trade/currencies"
	e "github.com/cryptounicorns/trade/errors"
	jsonTypes "github.com/cryptounicorns/trade/json"
	"github.com/cryptounicorns/trade/markets/market"
)

const (
	Name = "bitfinex"
	Addr = "https://api.bitfinex.com/v1"
)

var (
	Default       market.Market
	DefaultClient = http.DefaultClient
)

var (
	CurrencyMapping = map[currencies.Currency]string{
		// https://api.bitfinex.com/v1/symbols
		currencies.Bitcoin:            "btc",
		currencies.Litecoin:           "ltc",
		currencies.UnitedStatesDollar: "usd",
	}
	CurrencyPairDelimiter = ""
)

type Bitfinex struct {
	client *http.Client
}

type Ticker struct {
	High      jsonTypes.Float64String `json:"high"`
	Low       jsonTypes.Float64String `json:"low"`
	Vol       jsonTypes.Float64String `json:"volume"`
	Last      jsonTypes.Float64String `json:"last_price"`
	Buy       jsonTypes.Float64String `json:"bid"`
	Sell      jsonTypes.Float64String `json:"ask"`
	Timestamp jsonTypes.Float64String `json:"timestamp"`
}

//

func (m *Bitfinex) ID() string { return Name }

func (m *Bitfinex) GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	var (
		err error
	)

	tickers := make([]*market.Ticker, len(currencyPairs))
	for k, v := range currencyPairs {
		tickers[k], err = m.GetTicker(v)
		if err != nil {
			return nil, err
		}
	}

	return tickers, nil
}

func (m *Bitfinex) GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	var (
		u              *url.URL
		r              *http.Response
		pair           string
		ticker         *market.Ticker
		responseTicker = &Ticker{}
		err            error
	)

	u, err = url.Parse(Addr)
	if err != nil {
		return nil, err
	}

	u.Path += "/pubticker"

	pair, err = currencyPair.Format(
		CurrencyMapping,
		CurrencyPairDelimiter,
	)
	if err != nil {
		return nil, err
	}
	u.Path += "/" + pair

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
	err = json.NewDecoder(r.Body).Decode(responseTicker)
	if err != nil {
		return nil, err
	}

	// XXX: no avg and volcur fields
	ticker = market.NewTicker(m, currencyPair)
	ticker.Buy = float64(responseTicker.Buy)
	ticker.High = float64(responseTicker.High)
	ticker.Last = float64(responseTicker.Last)
	ticker.Low = float64(responseTicker.Low)
	ticker.Sell = float64(responseTicker.Sell)
	ticker.Timestamp = float64(responseTicker.Timestamp)
	ticker.Vol = float64(responseTicker.Vol)

	return ticker, nil
}

func (m *Bitfinex) Close() error { return nil }

//

func GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	return Default.GetTicker(currencyPair)
}

//

func New(c *http.Client) (*Bitfinex, error) {
	if c == nil {
		return nil, e.NewErrArgumentIsNil(c)
	}
	return &Bitfinex{c}, nil
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
