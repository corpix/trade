package dummy

import (
	"github.com/cryptounicorns/trade/markets/market"
)

func (m *Dummy) NewTickerConsumer() (market.TickerConsumer, error) {
	panic("Not implemented")
	return nil, nil
}

func NewTickerConsumer() (market.TickerConsumer, error) {
	return Default.TickerConsumer()
}
