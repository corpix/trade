package market

type Currency int

const (
	InvalidCurrency Currency = iota
	BTC
	LTC
	ETH
	GHS
	USD
	EUR
	RUB
)

var (
	CurrencyMapping = map[Currency]string{
		BTC: "BTC",
		LTC: "LTC",
		ETH: "ETH",
		GHS: "GHS",
		USD: "USD",
		EUR: "EUR",
		RUB: "RUB",
	}
	CurrencyPairDelimiter = "-"
)

func (c Currency) String() string {
	return CurrencyMapping[c]
}
