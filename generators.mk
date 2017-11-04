tools := ./tools

.PHONY: currencies/currencies.go
currencies/currencies.go:
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
	} | $(tools)/generate-currencies --package-name=currencies all > $@
	go fmt $@

.PHONY: markets/market/bitfinex/currency_mapping.go
markets/market/bitfinex/currency_mapping.go:
	go run $(tools)/coinmarketcap/coinmarketcap.go  \
		exchanges --exchange=bitfinex           \
		| $(tools)/generate-currencies          \
			--package-name=bitfinex         \
			exchange                        \
		> $@
	go fmt $@

.PHONY: generate
generate:: currencies/currencies.go
generate:: markets/market/bitfinex/currency_mapping.go
