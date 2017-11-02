package market

import (
	"io"
)

type Market interface {
	ID() string

	Connect() (io.ReadWriteCloser, error)
	NewTickerConsumer(io.ReadWriter) TickerConsumer
}
