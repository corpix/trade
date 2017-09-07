package yobit

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/cryptounicorns/trade/currencies"
	e "github.com/cryptounicorns/trade/errors"
	"github.com/cryptounicorns/trade/markets/market"
)

const (
	Name = "yobit"
	Addr = "https://yobit.net/api/2"
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

type Yobit struct {
	client *http.Client
}

type Ticker struct {
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Avg        float64 `json:"avg"`
	Vol        float64 `json:"vol"`
	VolCur     float64 `json:"vol_cur"`
	Last       float64 `json:"last"`
	Buy        float64 `json:"buy"`
	Sell       float64 `json:"sell"`
	Updated    int64   `json:"updated"`
	ServerTime int64   `json:"server_time"`
}

//

func (m *Yobit) ID() string { return Name }

func (m *Yobit) GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
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

func (m *Yobit) GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	var (
		u              *url.URL
		r              *http.Response
		pair           string
		ticker         *market.Ticker
		responseTicker = make(map[string]Ticker, 1)
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

	u.Path += "/" + pair + "/ticker"

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
	err = json.NewDecoder(r.Body).Decode(&responseTicker)
	if err != nil {
		return nil, err
	}

	ticker = market.NewTicker(m, currencyPair)
	ticker.Avg = responseTicker["ticker"].Avg
	ticker.Buy = responseTicker["ticker"].Buy
	ticker.High = responseTicker["ticker"].High
	ticker.Last = responseTicker["ticker"].Last
	ticker.Low = responseTicker["ticker"].Low
	ticker.Sell = responseTicker["ticker"].Sell
	ticker.Timestamp = float64(responseTicker["ticker"].Updated)
	ticker.Vol = responseTicker["ticker"].Vol
	ticker.VolCur = responseTicker["ticker"].VolCur

	return ticker, nil
}

func (m *Yobit) Close() error { return nil }

//

func GetTickers(currencyPairs []currencies.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair currencies.CurrencyPair) (*market.Ticker, error) {
	return Default.GetTicker(currencyPair)
}

//

func New(c *http.Client) (*Yobit, error) {
	if c == nil {
		return nil, e.NewErrArgumentIsNil(c)
	}
	return &Yobit{c}, nil
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
