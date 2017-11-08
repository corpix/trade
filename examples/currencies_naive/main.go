package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/cryptounicorns/trade/assets"
	"github.com/cryptounicorns/trade/currencies"
)

func main() {
	var (
		commonCurrencies []byte
		common           currencies.Currencies

		marketCurrencies []byte
		market           currencies.Currencies

		mapper currencies.Mapper

		commonDash currencies.Currency
		marketDash currencies.Currency

		err error
	)

	commonCurrencies, err = assets.Asset("currencies/currencies.json")
	if err != nil {
		panic(err)
	}

	marketCurrencies, err = assets.Asset("markets/market/bitfinex/currencies.json")
	if err != nil {
		panic(err)
	}

	common, err = currencies.NewFromJSON(commonCurrencies)
	if err != nil {
		panic(err)
	}

	market, err = currencies.NewFromJSON(marketCurrencies)
	if err != nil {
		panic(err)
	}

	mapper = currencies.NewMapper(
		common,
		market,
	)

	marketDash, err = mapper.ToCommon(market.MapByName()["dash"])
	if err != nil {
		panic(err)
	}

	commonDash, err = mapper.ToMarket(common.MapByName()["dash"])
	if err != nil {
		panic(err)
	}

	spew.Dump(
		commonDash,
		marketDash,
	)
}
