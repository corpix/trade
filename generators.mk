.PHONY: currencies/currencies.go
currencies/currencies.go:
	./currencies/currencies.py > $@
