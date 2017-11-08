package bitfinex

type symbolPairByChannel map[uint]SymbolPair

func (c symbolPairByChannel) Get(channelID uint) (SymbolPair, error) {
	var (
		pair SymbolPair
		ok   bool
	)

	pair, ok = c[channelID]
	if !ok {
		return pair, NewErrUnknownChannel(channelID)
	}

	return pair, nil
}
