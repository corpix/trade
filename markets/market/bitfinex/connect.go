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
		res ws.Response
		err error
	)

	ctx, cancelCtx = context.WithTimeout(
		context.Background(),
		m.config.Endpoint.Timeout,
	)
	defer cancelCtx()

	r, res, err = ws.DefaultDialer.Dial(
		ctx,
		m.config.Endpoint.URL.String(),
		m.config.Endpoint.Headers,
	)
	if err != nil {
		return nil, err
	}
	res.Body.Close()

	return r, nil
}
