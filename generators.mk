tools := ./tools

.PHONY: currencies/currencies.json
currencies/currencies.json:
	# FIXME: Fiat currencies are also here, but they should be delivered from some resource
	{                                                                                    \
		set -e;                                                                      \
		echo '{"name": "China Yan",            "symbol": "CNY", "volume": 9999999}'; \
		echo '{"name": "Japanese Yen",         "symbol": "JPY", "volume": 9999999}'; \
		echo '{"name": "Russian Ruble",        "symbol": "RUB", "volume": 9999999}'; \
		echo '{"name": "United States Dollar", "symbol": "USD", "volume": 9999999}'; \
		echo '{"name": "Euro",                 "symbol": "EUR", "volume": 9999999}'; \
		echo '{"name": "Canadian Dollar",      "symbol": "CAD", "volume": 9999999}'; \
		go run $(tools)/coinmarketcap/coinmarketcap.go all;                          \
	} | $(tools)/postprocess-currencies > $@

.PHONY: markets/market/bitfinex/currencies.json
markets/market/bitfinex/currencies.json:
	go run $(tools)/coinmarketcap/coinmarketcap.go  \
		exchanges --exchange=bitfinex           \
		| $(tools)/postprocess-currencies       \
		> $@

.PHONY: generate
generate:: currencies/currencies.json
generate:: markets/market/bitfinex/currencies.json
