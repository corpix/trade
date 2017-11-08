package currencies

type CurrencyPair struct {
	Left  Currency
	Right Currency
}

func NewCurrencyPair(left Currency, right Currency) CurrencyPair {
	return CurrencyPair{
		Left:  left,
		Right: right,
	}
}
