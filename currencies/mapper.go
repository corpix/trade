package currencies

type Mapper struct {
	commonToMarket Mapping
	commonByName   map[string]Currency
	commonBySymbol map[Symbol]Currency

	marketToCommon Mapping
	marketByName   map[string]Currency
	marketBySymbol map[Symbol]Currency
}

func (m Mapper) ToMarket(common Currency) (Currency, error) {
	var (
		c  Currency
		ok bool
	)

	c, ok = m.commonToMarket[common]
	if !ok {
		return c, NewErrNoMapping(
			common,
			MapperDirection{
				From: MapperDirectionCommon,
				To:   MapperDirectionMarket,
			},
		)
	}

	return c, nil
}

func (m Mapper) CommonByName(name string) (Currency, error) {
	var (
		c  Currency
		ok bool
	)

	c, ok = m.commonByName[name]
	if !ok {
		return c, NewErrNotFound(name)
	}

	return c, nil
}

func (m Mapper) CommonBySymbol(symbol Symbol) (Currency, error) {
	var (
		c  Currency
		ok bool
	)

	c, ok = m.commonBySymbol[symbol]
	if !ok {
		return c, NewErrSymbolNotFound(symbol)
	}

	return c, nil
}

func (m Mapper) ToCommon(market Currency) (Currency, error) {
	var (
		c  Currency
		ok bool
	)

	c, ok = m.marketToCommon[market]
	if !ok {
		return c, NewErrNoMapping(
			market,
			MapperDirection{
				From: MapperDirectionMarket,
				To:   MapperDirectionCommon,
			},
		)
	}

	return c, nil
}

func (m Mapper) MarketByName(name string) (Currency, error) {
	var (
		c  Currency
		ok bool
	)

	c, ok = m.marketByName[name]
	if !ok {
		return c, NewErrNotFound(name)
	}

	return c, nil
}

func (m Mapper) MarketBySymbol(symbol Symbol) (Currency, error) {
	var (
		c  Currency
		ok bool
	)

	c, ok = m.marketBySymbol[symbol]
	if !ok {
		return c, NewErrSymbolNotFound(symbol)
	}

	return c, nil
}

func NewMapper(common Currencies, market Currencies) Mapper {
	return Mapper{
		commonToMarket: NewMappingFromIntersection(
			common,
			market,
		),
		commonByName:   common.MapByName(),
		commonBySymbol: common.MapBySymbol(),

		marketToCommon: NewMappingFromIntersection(
			market,
			common,
		),
		marketByName:   market.MapByName(),
		marketBySymbol: market.MapBySymbol(),
	}
}
