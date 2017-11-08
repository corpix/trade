package market

import (
	"io"
)

type Market interface {
	Name() string
	Connect() (io.ReadWriteCloser, error)
	NewTickerConsumer(io.ReadWriter) TickerConsumer
}
