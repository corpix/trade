package btce

// The MIT License (MIT)
//
// Copyright © 2017 Dmitry Moskowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
	DefaultClient = http.DefaultClient
	Default       market.Market
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
		pair, err = market.CurrencyPairFromString(
			k,
			CurrencyMapping,
			CurrencyPairDelimiter,
		)
		if err != nil {
			return nil, err
		}

		tickers[n] = market.NewTicker(m, pair)
		tickers[n].Timestamp = float64(v.Updated)
		err = copier.Copy(&tickers[n], v)
		if err != nil {
			return nil, err
		}
		n++
	}

	return tickers, nil
}

func (m *Btce) GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
	tickers, err := m.GetTickers(
		[]market.CurrencyPair{currencyPair},
	)
	if err != nil {
		return nil, err
	}

	return tickers[0], nil
}

func (m *Btce) Close() error { return nil }

//

func GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
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
