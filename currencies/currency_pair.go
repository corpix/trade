package currencies

type currencyPairJSON [2]Currency

type CurrencyPair struct {
	Left  Currency
	Right Currency
}

func (c CurrencyPair) MarshalJSON() ([]byte, error) {
	return jsonFormat.Marshal(
		currencyPairJSON{
			c.Left,
			c.Right,
		},
	)
}

func (c CurrencyPair) UnmarshalJSON(buf []byte) error {
	var (
		pair = currencyPairJSON{}
		err  error
	)

	err = jsonFormat.Unmarshal(buf, pair)
	if err != nil {
		return err
	}

	c.Left = pair[0]
	c.Right = pair[1]

	return nil
}

func NewCurrencyPair(left Currency, right Currency) CurrencyPair {
	return CurrencyPair{
		Left:  left,
		Right: right,
	}
}
