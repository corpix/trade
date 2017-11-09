package bitfinex

import (
	"context"
	"io"

	"github.com/gobwas/ws"
)

var (
	ctx       context.Context
	cancelCtx context.CancelFunc
)

func (m *Bitfinex) Connect() (io.ReadWriteCloser, error) {
	var (
		r   io.ReadWriteCloser
		err error
	)

	ctx, cancelCtx = context.WithTimeout(
		context.Background(),
		m.config.Endpoint.Timeout,
	)
	defer cancelCtx()

	r, _, _, err = ws.DefaultDialer.Dial(
		ctx,
		m.config.Endpoint.URL.String(),
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}
