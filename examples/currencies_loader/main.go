package main

import (
	"fmt"

	"github.com/cryptounicorns/trade/assets"
	"github.com/cryptounicorns/trade/currencies"
)

func main() {
	var (
		loader = currencies.NewAssetLoader(assets.Asset)

		marketName           = "bitfinex"
		currencyName         = "dash"
		marketCurrencySymbol = currencies.Symbol("DSH")
		commonCurrencySymbol = currencies.Symbol("DASH")

		common currencies.Currencies
		market currencies.Currencies

		mapper currencies.Mapper

		commonCurrency currencies.Currency
		marketCurrency currencies.Currency

		err error
	)

	common, err = loader.Common()
	if err != nil {
		panic(err)
	}
	market, err = loader.Market(marketName)
	if err != nil {
		panic(err)
	}

	mapper = currencies.NewMapper(
		common,
		market,
	)

	// the following code could be viewed as 3 different examples
	// for different use-cases ->

	// mapped by currency name
	commonCurrency, err = mapper.CommonByName(currencyName)
	if err != nil {
		panic(err)
	}
	marketCurrency, err = mapper.MarketByName(currencyName)
	if err != nil {
		panic(err)
	}

	fmt.Println("Mapping by currency name:")
	fmt.Printf(
		"On %s market currency %s has symbol %s\nBut in common we use symbol %s\n",
		marketName,
		currencyName,
		marketCurrency.Symbol,
		commonCurrency.Symbol,
	)
	fmt.Println("---")

	// mapped by currency symbol
	commonCurrency, err = mapper.CommonBySymbol(commonCurrencySymbol)
	if err != nil {
		panic(err)
	}
	marketCurrency, err = mapper.MarketBySymbol(marketCurrencySymbol)
	if err != nil {
		panic(err)
	}

	fmt.Println("Mapping by currency symbol:")
	fmt.Printf(
		"On %s market currency %s has symbol %s\nBut in common we use symbol %s\n",
		marketName,
		currencyName,
		marketCurrency.Symbol,
		commonCurrency.Symbol,
	)
	fmt.Println("---")

	// mapped by currency structure
	commonCurrency, err = mapper.ToCommon(marketCurrency)
	if err != nil {
		panic(err)
	}
	marketCurrency, err = mapper.ToMarket(commonCurrency)
	if err != nil {
		panic(err)
	}

	fmt.Println("Mapping by currency structure:")
	fmt.Printf(
		"On %s market currency %s has symbol %s\nBut in common we use symbol %s\n",
		marketName,
		currencyName,
		marketCurrency.Symbol,
		commonCurrency.Symbol,
	)
	fmt.Println("---")
}
