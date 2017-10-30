package bitfinex

import (
	"net/http"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/logrus"
	"github.com/corpix/loggers/logger/prefixwrapper"
	logrusLogger "github.com/sirupsen/logrus"

	e "github.com/cryptounicorns/trade/errors"
	"github.com/cryptounicorns/trade/markets/market"
)

const (
	Name = "bitfinex"
	Addr = "https://api.bitfinex.com/v1"
)

var (
	Default       market.Market
	DefaultClient = http.DefaultClient
)

func init() {
	var (
		err error
	)

	Default, err = New(
		DefaultClient,
		Config{},
		logrus.New(logrusLogger.New()),
	)
	if err != nil {
		panic(err)
	}
}

type Bitfinex struct {
	client *http.Client
	config Config
	log    loggers.Logger
}

func (m *Bitfinex) ID() string { return Name }

func (m *Bitfinex) Close() error { return nil }

func New(c *http.Client, cf Config, l loggers.Logger) (*Bitfinex, error) {
	if c == nil {
		return nil, e.NewErrArgumentIsNil(c)
	}
	return &Bitfinex{
		client: c,
		config: cf,
		log: prefixwrapper.New(
			Name+": ",
			l,
		),
	}, nil
}
