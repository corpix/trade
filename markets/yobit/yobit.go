package yobit

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

	"encoding/json"
	e "github.com/corpix/trade/errors"
	"github.com/corpix/trade/market"
	"net/url"
)

const (
	Name = "yobit"
	Addr = "https://yobit.net/api/2"
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

func (m *Yobit) GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
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

func (m *Yobit) GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
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

func (m *Yobit) Close() error { return nil }

//

func GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
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
