package bitfinex

import (
	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"

	"github.com/cryptounicorns/trade/currencies"
)

const (
	Name    = "bitfinex"
	Version = 2
)

type Bitfinex struct {
	config     Config
	currencies currencies.Mapper
	log        loggers.Logger
}

func (m *Bitfinex) Name() string {
	return Name
}

func New(config Config, mapper currencies.Mapper, log loggers.Logger) *Bitfinex {
	return &Bitfinex{
		config:     config,
		currencies: mapper,
		log: prefixwrapper.New(
			Name+": ",
			log,
		),
	}
}
