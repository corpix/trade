package markets

import (
	"net/http"
	"strings"

	"github.com/corpix/trade/markets/market"
	"github.com/corpix/trade/markets/market/bitfinex"
	"github.com/corpix/trade/markets/market/btce"
	"github.com/corpix/trade/markets/market/cex"
	"github.com/corpix/trade/markets/market/yobit"
)

const (
	BitfinexMarket = bitfinex.Name
	BtceMarket     = btce.Name
	CexMarket      = cex.Name
	YobitMarket    = yobit.Name
)

var (
	DefaultMarkets = map[string]market.Market{
		BitfinexMarket: bitfinex.Default,
		BtceMarket:     btce.Default,
		CexMarket:      cex.Default,
		YobitMarket:    yobit.Default,
	}

	DefaultClients = map[string]interface{}{
		BitfinexMarket: bitfinex.DefaultClient,
		BtceMarket:     btce.DefaultClient,
		CexMarket:      cex.DefaultClient,
		YobitMarket:    yobit.DefaultClient,
	}
)

func GetDefault(market string) (market.Market, error) {
	if m, ok := DefaultMarkets[strings.ToLower(market)]; ok {
		return m, nil
	}
	return nil, NewErrUnsupportedMarket(market)
}

func GetDefaultClient(market string) (interface{}, error) {
	if t, ok := DefaultClients[strings.ToLower(market)]; ok {
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

	// FIXME: Type assertion for clients look like crap
	// Probably this should be concrete type.
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
