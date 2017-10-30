package dummy

import (
	"net/http"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"

	e "github.com/cryptounicorns/trade/errors"
	"github.com/cryptounicorns/trade/markets/market"
)

const (
	Name = "dummy"
	Addr = "https://localhost"
)

var (
	Default       market.Market
	DefaultClient = http.DefaultClient
)

func init() {
	var (
		err error
	)

	Default, err = New(DefaultClient)
	if err != nil {
		panic(err)
	}
}

type Dummy struct {
	client *http.Client
	log    loggers.Logger
}

func (m *Dummy) ID() string { return Name }

func (m *Dummy) Close() error { return nil }

func New(c *http.Client, l loggers.Logger) (*Dummy, error) {
	if c == nil {
		return nil, e.NewErrArgumentIsNil(c)
	}
	if l == nil {
		return nil, e.NewErrArgumentIsNil(l)
	}

	return &Dummy{
		client: c,
		log: prefixwrapper.New(
			Name+": ",
			l,
		),
	}, nil
}
