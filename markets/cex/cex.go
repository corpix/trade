package cex

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
	High      jsonTypes.Float64String `json:"high"`
	Low       jsonTypes.Float64String `json:"low"`
	Vol       jsonTypes.Float64String `json:"volume30d"`
	VolCur    jsonTypes.Float64String `json:"volume"`
	Last      jsonTypes.Float64String `json:"last"`
	Buy       float64                 `json:"bid"`
	Sell      float64                 `json:"ask"`
	Timestamp jsonTypes.Int64String   `json:"timestamp"`
	Error     string                  `json:"error"`
}

//

func (m *Cex) ID() string { return Name }

func (m *Cex) GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
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

	pair, err = currencyPair.Format(
		CurrencyMapping,
		CurrencyPairDelimiter,
	)
	if err != nil {
		return nil, err
	}

	u.Path += "/ticker/" + pair

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

	if len(responseTicker.Error) > 0 {
		return nil, NewErrApi(responseTicker.Error)
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

func (m *Cex) GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
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
