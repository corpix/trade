package bitfinex

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/jinzhu/copier"

	e "github.com/corpix/trade/errors"
	jsonTypes "github.com/corpix/trade/json"
	"github.com/corpix/trade/market"
)

const (
	Addr = "https://api.bitfinex.com/v1"
)

var (
	DefaultTransport *Transport
	Default          market.Market
)

var (
	CurrencyMapping = map[market.Currency]string{
		// https://api.bitfinex.com/v1/symbols
		market.BTC: "btc",
		market.LTC: "ltc",
		market.USD: "usd",
	}
	CurrencyPairDelimiter = ""
)

type Transport struct {
	Addr   string
	Client *http.Client
}

type Bitfinex struct {
	transport *Transport
}

type Ticker struct {
	High jsonTypes.Float64String `json:"high"`
	Low  jsonTypes.Float64String `json:"low"`
	// FIXME: Bitfinex does not have avg :\
	// Yare yare
	Avg jsonTypes.Float64String `json:"avg"`
	Vol jsonTypes.Float64String `json:"volume"`
	// FIXME: Bitfinex does not have volcur :\
	// Yare yare
	VolCur    jsonTypes.Float64String `json:"vol_cur"`
	Last      jsonTypes.Float64String `json:"last_price"`
	Buy       jsonTypes.Float64String `json:"bid"`
	Sell      jsonTypes.Float64String `json:"ask"`
	Timestamp jsonTypes.Float64String `json:"timestamp"`
}

//

func (m *Bitfinex) ID() string { return "bitfinex" }

func (m *Bitfinex) GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
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

func (m *Bitfinex) GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
	var (
		u              *url.URL
		r              *http.Response
		pair           string
		ticker         *market.Ticker
		responseTicker = &Ticker{}
		err            error
	)

	u, err = url.Parse(m.transport.Addr)
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

	r, err = m.transport.Client.Get(u.String())
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

	ticker = market.NewTicker(m, currencyPair)
	err = copier.Copy(&ticker, responseTicker)
	if err != nil {
		return nil, err
	}

	return ticker, nil
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

func New(transport *Transport) (*Bitfinex, error) {
	if transport == nil {
		return nil, e.NewErrArgumentIsNil(transport)
	}
	return &Bitfinex{transport}, nil
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
