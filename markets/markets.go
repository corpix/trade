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
	"sync"

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
		tickers = []*market.Ticker{}
		w       = &sync.WaitGroup{}
		ts      = make(chan *market.Ticker)
		errs    = make(chan error)
	)
	defer close(errs)

	for _, v := range markets {
		w.Add(1)
		m.Pool.Feed <- pool.NewWork(
			context.TODO(),
			func(ctx context.Context) {
				buf, err := v.GetTickers(currencyPairs)
				// FIXME: This part could be a custom work type in pool package
				// (with error handling)
				if err != nil {
					errs <- err
					return
				}
				for _, v := range buf {
					ts <- v
				}
				w.Done()
			},
		)
	}

loop:
	for {
		select {
		case err := <-errs:
			// FIXME: Cancel works
			return nil, err
		case ticker := <-ts:
			tickers = append(
				tickers,
				ticker,
			)
			if len(tickers) == len(markets) {
				break loop
			}
		}
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
