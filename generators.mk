tools_root         := ./tools
assets_root        := ./assets
markets_root       := ./markets/market
currencies_root    := ./currencies
markets            := $(shell find $(markets_root)/* -maxdepth 1 -type d | xargs basename)

all_currencies     := $(currencies_root)/currencies.json
markets_currencies := $(foreach market,$(markets),$(markets_root)/$(market)/currencies.json)

.PHONY: $(all_currencies)
$(all_currencies):
	# FIXME: Fiat should be automatically downloaded for each market!
	{                                                                                    \
		set -e;                                                                      \
		echo '{"name": "china-yan",            "symbol": "CNY", "volume": 9999999}'; \
		echo '{"name": "japanese-yen",         "symbol": "JPY", "volume": 9999999}'; \
		echo '{"name": "russian-ruble",        "symbol": "RUB", "volume": 9999999}'; \
		echo '{"name": "united-states-dollar", "symbol": "USD", "volume": 9999999}'; \
		echo '{"name": "euro",                 "symbol": "EUR", "volume": 9999999}'; \
		echo '{"name": "canadian-dollar",      "symbol": "CAD", "volume": 9999999}'; \
		go run $(tools_root)/coinmarketcap/coinmarketcap.go all;                     \
	} | $(tools_root)/postprocess-currencies --verbose > $@

.PHONY: $(markets_currencies)
$(markets_currencies):
	# FIXME: Fiat should be automatically downloaded for each market!
	{                                                                                    \
		set -e;                                                                      \
		echo '{"name": "united-states-dollar", "symbol": "USD", "volume": 9999999}'; \
		go run $(tools_root)/coinmarketcap/coinmarketcap.go                          \
			exchanges --exchange=$(shell basename $(@:%/currencies.json=%));     \
	} | $(tools_root)/postprocess-currencies --verbose > $@

.PHONY: generate
generate:: $(all_currencies)
generate:: $(markets_currencies)
generate::
	go-bindata                                                      \
		-o        $(assets_root)/assets.go                      \
		--pkg     assets                                        \
		--ignore '.*\.go'                                       \
		--nometadata                                            \
		$(foreach market,$(markets),$(markets_root)/$(market)/) \
		$(currencies_root)/
	go fmt $(assets_root)/assets.go
