package bitfinex

import (
	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/logrus"
	"github.com/corpix/loggers/logger/prefixwrapper"
	logrusLogger "github.com/sirupsen/logrus"
)

const (
	Name    = "bitfinex"
	Version = 2
)

var (
	Default       = New(DefaultConfig, DefaultLogger)
	DefaultLogger = logrus.New(logrusLogger.New())
)

type Bitfinex struct {
	config Config
	log    loggers.Logger
}

func (m *Bitfinex) ID() string { return Name }

func New(cf Config, l loggers.Logger) *Bitfinex {
	return &Bitfinex{
		config: cf,
		log: prefixwrapper.New(
			Name+": ",
			l,
		),
	}
}
