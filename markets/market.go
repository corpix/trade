package markets

import (
	"io"

	"github.com/cryptounicorns/trade/markets/ticker"
)

type Market interface {
	Name() string
	Connect() (io.ReadWriteCloser, error)
	NewTickerConsumer(io.ReadWriter) (ticker.Consumer, error)
}
