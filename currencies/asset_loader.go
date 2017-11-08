package currencies

type AssetLoader func(name string) ([]byte, error)

func (l AssetLoader) Common() (Currencies, error) {
	var (
		c   []byte
		err error
	)

	c, err = l("currencies/currencies.json")
	if err != nil {
		return nil, err
	}

	return NewFromJSON(c)
}

func (l AssetLoader) Market(market string) (Currencies, error) {
	var (
		m   []byte
		err error
	)

	m, err = l("markets/market/" + market + "/currencies.json")
	if err != nil {
		return nil, err
	}

	return NewFromJSON(m)
}

func NewAssetLoader(fn func(name string) ([]byte, error)) Loader {
	return AssetLoader(fn)
}
