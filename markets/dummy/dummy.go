package dummy

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

	e "github.com/corpix/trade/errors"
	"github.com/corpix/trade/market"
	transport "github.com/corpix/trade/transport/http"
)

const (
	Addr = "https://localhost"
)

var (
	DefaultTransport *transport.Transport
	Default          market.Market
)

var (
	CurrencyMapping = map[market.Currency]string{
		market.BTC: "btc",
		market.LTC: "ltc",
		market.USD: "usd",
		market.EUR: "eur",
		market.RUB: "rub",
	}
	CurrencyPairDelimiter = "-"
)

type Dummy struct {
	transport *transport.Transport
}

//

func (m *Dummy) ID() string { return "dummy" }

func (m *Dummy) GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	tickers := make([]*market.Ticker, len(currencyPairs))
	for k, v := range currencyPairs {
		tickers[k] = market.NewTicker(
			m,
			v,
		)
	}
	return tickers, nil
}

func (m *Dummy) GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
	return market.NewTicker(
		m,
		currencyPair,
	), nil
}

//

func GetTickers(currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(currencyPairs)
}

func GetTicker(currencyPair market.CurrencyPair) (*market.Ticker, error) {
	return Default.GetTicker(currencyPair)
}

//

func New(t *transport.Transport) (*Dummy, error) {
	if t == nil {
		return nil, e.NewErrArgumentIsNil(t)
	}
	return &Dummy{t}, nil
}

//

func init() {
	var (
		err error
	)

	DefaultTransport, err = transport.New(
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
