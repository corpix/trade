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
	"net/http"
	"strings"

	"github.com/corpix/trade/market"
	"github.com/corpix/trade/markets/bitfinex"
	"github.com/corpix/trade/markets/btce"
	"github.com/corpix/trade/markets/cex"
	"github.com/corpix/trade/markets/yobit"
)

const (
	BitfinexMarket = bitfinex.Name
	BtceMarket     = btce.Name
	CexMarket      = cex.Name
	YobitMarket    = yobit.Name
)

var (
	Markets = map[string]market.Market{
		BitfinexMarket: bitfinex.Default,
		BtceMarket:     btce.Default,
		CexMarket:      cex.Default,
		YobitMarket:    yobit.Default,
	}

	Clients = map[string]interface{}{
		BitfinexMarket: bitfinex.DefaultClient,
		BtceMarket:     btce.DefaultClient,
		CexMarket:      cex.DefaultClient,
		YobitMarket:    yobit.DefaultClient,
	}
)

func GetDefault(market string) (market.Market, error) {
	if m, ok := Markets[strings.ToLower(market)]; ok {
		return m, nil
	}
	return nil, NewErrUnsupportedMarket(market)
}

func GetDefaultClient(market string) (interface{}, error) {
	if t, ok := Clients[strings.ToLower(market)]; ok {
		return t, nil
	}
	return nil, NewErrUnsupportedMarket(market)
}

func New(market string, client interface{}) (market.Market, error) {
	switch client.(type) {
	case *http.Client:
	default:
		return nil, NewErrUnsupportedClient(client)
	}

	switch strings.ToLower(market) {
	case BitfinexMarket:
		return bitfinex.New(client.(*http.Client))
	case BtceMarket:
		return btce.New(client.(*http.Client))
	case CexMarket:
		return cex.New(client.(*http.Client))
	case YobitMarket:
		return yobit.New(client.(*http.Client))
	default:
		return nil, NewErrUnsupportedMarket(market)
	}
}
