## Currency and market names

### Download

> There are large amount of exchanges in the world. Some of them have different naming for
> the same assets(currencies). To make world more consistent we download this names from the
> single place which is a [coinmarketcap.com](https://coinmarketcap.com).

We have a tool in `tools/coinmarketcap` which scrapes some coinmarketcap pages and transforms the
data into JSON format so we could postprocess the data.

Here is a `--help` output of the tool:

``` console

  λ  go run ./tools/coinmarketcap/coinmarketcap.go --help

NAME:
   coinmarketcap - A new cli application

USAGE:
   coinmarketcap [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     all        Get data about all coins
     exchanges  Get data about concrete exchange coins
     help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --format value  Name of the output format supported by https://github.com/corpix/formats (default: "json")
   --help, -h      show help
   --version, -v   print the version
```

We could download data for all currencies:

``` console

  λ  go run ./tools/coinmarketcap/coinmarketcap.go all | head -n 3

{"name":"bitcoin","symbol":"BTC","volume":2722850000}
{"name":"ethereum","symbol":"ETH","volume":540705000}
{"name":"bitcoin-cash","symbol":"BCH","volume":1425000000}
...
```

Or download data for a single exchange:

``` console

  λ  go run ./tools/coinmarketcap/coinmarketcap.go exchanges --exchange=bitfinex | head -n 3

{"name":"bitcoin","symbol":"BTC","volume":352706000}
{"name":"bitcoin-cash","symbol":"BCH","volume":167144130}
{"name":"ethereum","symbol":"ETH","volume":85067700}
...
```

### Postprocess

Downloaded data could be postprocessed with `tools/postprocess-currencies`. What postprocessor does:

- Filters out currencies with 24h volume < 5000 $
- Transforms names to lowercase
- Transforms symbols to uppercase
- Removing information about volume(because it is not required in our meta-data)

Here is the `--help` of the tool:

``` console

  λ  ./tools/postprocess-currencies --help

usage: postprocess-currencies [-h] [--min-volume MIN_VOLUME] [--verbose]

optional arguments:
  -h, --help            show this help message and exit
  --min-volume MIN_VOLUME
                        Minimal volume in USD to be listed
  --verbose             Be more verbose
```

### Automation

Previous steps are automated as GNU make targets in `generators.mk`.

It serves a `generate` target. This target generates JSON files with:

- all currency names
- currency names used on particular exchange

### Caveats

- Fiat currencies are hardcoded at this time at `currencies/fiat_currencies.json` and should be checkined onto VCS
