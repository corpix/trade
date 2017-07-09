package markets

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
