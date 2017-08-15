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
