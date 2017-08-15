package currencies

type Currency int

func (c Currency) String() string {
	return CurrencyMapping[c]
}
