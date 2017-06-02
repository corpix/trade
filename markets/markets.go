package markets

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
	"context"

	"github.com/corpix/pool"

	"github.com/corpix/trade/market"
)

var (
	Default = New(50, 150)
)

type Markets struct {
	*pool.Pool
}

func (m *Markets) GetTickers(markets []market.Market, currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	var (
		works   = len(markets)
		tickers = []*market.Ticker{}
		results = make(chan *pool.Result)
	)
	defer close(results)

	for _, v := range markets {
		ctx := context.WithValue(
			context.Background(),
			"market",
			v,
		)
		m.Pool.Feed <- pool.NewWorkWithResult(
			ctx,
			results,
			func(ctx context.Context) (interface{}, error) {
				buf, err := ctx.
					Value("market").(market.Market).
					GetTickers(currencyPairs)
				if err != nil {
					return nil, err
				}
				return buf, nil
			},
		)
	}

	for works > 0 {
		result := <-results
		if result.Err != nil {
			return nil, result.Err
		}
		tickers = append(
			tickers,
			result.Value.([]*market.Ticker)...,
		)
		works--
	}

	return tickers, nil
}

//

func GetTickers(markets []market.Market, currencyPairs []market.CurrencyPair) ([]*market.Ticker, error) {
	return Default.GetTickers(markets, currencyPairs)
}

//

func New(workers, queueSize int) *Markets {
	return &Markets{pool.New(workers, queueSize)}
}
