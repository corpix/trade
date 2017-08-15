trade
-----

[![Build Status](https://travis-ci.org/corpix/trade.svg?branch=master)](https://travis-ci.org/corpix/trade)

Consistent crypto currency markets trading api client.

## Currency names

At this time all currency names are loaded from [coinmarketcap](https://coinmarketcap.com/all/views/all/).

There is a script [currencies/currencies.py](currencies/currencies.py) which should be called with make:

``` console
$ make currencies/currencies.go
```

This script will download all cryptocurrency data from coinmarketcap and generate the Go code from this data.

There are some caveats because of automation nature of the process:

- If currency name begins with numbers(such as `1337`) then it will be prepended with `Coin`
- Not all currencies will be present, only currencies which have non `low cap` on coinmarketcap
- I don't know where could I get fiat currencies names so at this time fiat currencies list is hardcoded

Please look [currencies/currencies.go](currencies/currencies.go) for generated Go source.
