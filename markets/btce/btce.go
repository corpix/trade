package btce

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/jinzhu/copier"

	e "github.com/corpix/trade/errors"
	"github.com/corpix/trade/market"
)

const (
	Addr = "https://btc-e.com/api/3"
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
		market.RUB: "rur",
	}
	CurrencyPairDelimiter = "_"
)

type Transport struct {
	Addr   string
	Client *http.Client
}

type Btce struct {
	transport *Transport
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

func (m *Btce) ID() string { return "btce" }

func (m *Btce) GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	var (
		u               *url.URL
		r               *http.Response
		n               int
		pair            market.CurrencyPair
		pairs           = make([]string, len(currencyPairs))
		responseTickers = make(map[string]Ticker, len(currencyPairs))
		tickers         = make([]*market.Ticker, len(currencyPairs))
		err             error
	)

	u, err = url.Parse(m.transport.Addr)
	if err != nil {
		return nil, err
	}

	u.Path += "/ticker"

	for k, v := range currencyPairs {
		pairs[k], err = m.CurrencyPair(v)
		if err != nil {
			return nil, err
		}
	}
	u.Path += "/" + strings.Join(
		pairs,
		"-",
	)

	//

	r, err = m.transport.Client.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&responseTickers)
	if err != nil {
		return nil, err
	}

	n = 0
	for k, v := range responseTickers {
		pair, err = market.CurrencyPairFromString(
			k,
			CurrencyMapping,
			CurrencyPairDelimiter,
		)
		if err != nil {
			return nil, err
		}

		tickers[n] = market.NewTicker(m, pair)
		err = copier.Copy(&tickers[n], v)
		if err != nil {
			return nil, err
		}
		n++
	}

	return tickers, nil
}

func (m *Btce) GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
	return nil, nil
}

func (m *Btce) CurrencyPair(currencyPair market.CurrencyPair) (string, error) {
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

func New(transport *Transport) (*Btce, error) {
	if transport == nil {
		return nil, e.NewErrArgumentIsNil(transport)
	}
	return &Btce{transport}, nil
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
