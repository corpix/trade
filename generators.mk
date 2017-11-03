tools := ./tools

.PHONY: currencies/currencies.go
currencies/currencies.go:
	go run    $(tools)/coinmarketcap/coinmarketcap.go all \
		| $(tools)/generate-currencies --verbose      \
		> $@
	go fmt $@
