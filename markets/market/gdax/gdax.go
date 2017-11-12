package gdax

import (
	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"

	"github.com/cryptounicorns/trade/currencies"
)

const (
	Name    = "gdax"
	Version = 2
)

type Gdax struct {
	config     Config
	currencies currencies.Mapper
	log        loggers.Logger
}

func (m *Gdax) Name() string {
	return Name
}

func New(config Config, mapper currencies.Mapper, log loggers.Logger) *Gdax {
	return &Gdax{
		config:     config,
		currencies: mapper,
		log: prefixwrapper.New(
			Name+": ",
			log,
		),
	}
}
