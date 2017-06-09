package bitfinex

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

	"github.com/jinzhu/copier"

	e "github.com/corpix/trade/errors"
	jsonTypes "github.com/corpix/trade/json"
	"github.com/corpix/trade/market"
)

const (
	Name = "bitfinex"
	Addr = "https://api.bitfinex.com/v1"
)

var (
	DefaultClient = http.DefaultClient
	Default       market.Market
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

type Bitfinex struct {
	client *http.Client
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

func (m *Bitfinex) ID() string { return Name }

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

	ticker = market.NewTicker(m, currencyPair)
	err = copier.Copy(&ticker, responseTicker)
	if err != nil {
		return nil, err
	}

	return ticker, nil
}

func (m *Bitfinex) Close() error { return nil }

//

func GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
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
