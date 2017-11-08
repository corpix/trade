package bitfinex

import (
	"io"
	"strconv"
	"time"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/cryptounicorns/websocket/consumer"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/market"
)

const (
	TickerChannelName = "ticker"
)

type TickerConsumer struct {
	*consumer.Consumer

	connection          io.ReadWriter
	currencies          currencies.Mapper
	channelToSymbolPair symbolPairByChannel
	tickers             chan *market.Ticker
	done                chan struct{}
	log                 loggers.Logger
}

func (c *TickerConsumer) subscribe(pair SymbolPair, iterator *Iterator) (uint, error) {
	var (
		event = Event{
			Event: SubscribeEventName,
		}
		e   []byte
		err error
	)

	e, err = Format.Marshal(
		&SubscribeTickerEvent{
			SubscribeEvent: SubscribeEvent{
				Event:   event,
				Channel: TickerChannelName,
			},
			Pair: pair,
		},
	)
	if err != nil {
		return 0, err
	}

	err = wsutil.WriteClientText(
		c.connection,
		e,
	)
	if err != nil {
		return 0, err
	}

	e, err = iterator.NextEvent()
	if err != nil {
		return 0, err
	}

	err = Format.Unmarshal(
		e,
		&event,
	)
	if err != nil {
		return 0, err
	}

	switch event.Event {
	case SubscribedEventName:
		subscribedEvent := &SubscribedEvent{}
		err = Format.Unmarshal(
			e,
			subscribedEvent,
		)
		if err != nil {
			return 0, err
		}

		return subscribedEvent.ChanID, nil
	case ErrorEventName:
		errorEvent := &ErrorEvent{}
		err = Format.Unmarshal(
			e,
			errorEvent,
		)
		if err != nil {
			return 0, err
		}

		return 0, NewErrSubscription(
			errorEvent.Channel,
			errorEvent.Msg,
		)
	default:
		return 0, NewErrUnexpectedEvent(
			SubscribeEventName+"|"+ErrorEventName,
			event.Event,
		)
	}
}

func (c *TickerConsumer) preamble(pairs []SymbolPair, iterator *Iterator) error {
	var (
		handshaker = NewHandshaker(iterator, c.log)

		channelID uint
		err       error
	)

	err = handshaker.Handshake()
	if err != nil {
		return err
	}
	c.log.Debug("Handshaked")

	for _, pair := range pairs {
		channelID, err = c.subscribe(pair, iterator)
		if err != nil {
			return err
		}

		c.channelToSymbolPair[channelID] = pair
		c.log.Debug("Subscribed ", channelID, pair)
	}

	c.log.Debugf(
		"Preamble complete, channels subscribed: %#v",
		c.channelToSymbolPair,
	)

	return nil
}

func (c *TickerConsumer) consume(iterator *Iterator) (*pairTicker, error) {
	var (
		expectedLen = 2
		data        = make(
			Data,
			expectedLen,
		)

		d []byte

		ticker    = Ticker{}
		channelID int

		pair SymbolPair
		err  error
	)

	d, err = iterator.NextData()
	if err != nil {
		return nil, err
	}

	err = Format.Unmarshal(d, &data)
	if err != nil {
		return nil, err
	}

	if len(data) != expectedLen {
		return nil, NewErrDataLengthMismatch(
			expectedLen,
			len(data),
		)
	}

	if len(data[1]) == 0 {
		return nil, NewErrEmptyDataPayload()
	}

	switch data[1][0] {
	case '[':
		// We got ticker, this is what we have expect.
	case '"':
		// We got string message(heartbeat), nothing to do with them
		// now, skipping.
		return nil, errContinue
	default:
		// FIXME: I don't like this error message
		// both arguments should represent type
		// but it is hard to infer it from string
		return nil, NewErrUnexpectedDataPayloadType(
			"Ticker",
			string(data[1]),
		)
	}

	err = Format.Unmarshal(data[1], &ticker)
	if err != nil {
		return nil, err
	}

	channelID, err = strconv.Atoi(
		string(data[0]),
	)
	if err != nil {
		return nil, err
	}

	pair, err = c.channelToSymbolPair.Get(
		uint(channelID),
	)
	if err != nil {
		return nil, err
	}

	return &pairTicker{
		SymbolPair: pair,
		Ticker:     ticker,
	}, nil
}

func (c *TickerConsumer) convertTicker(pt *pairTicker) *market.Ticker {
	// see: https://docs.bitfinex.com/v2/reference#ws-public-ticker
	// (snapshot)
	// [
	// 	CHANNEL_ID,
	// 	[
	// 	0	BID,
	// 	1	BID_SIZE,
	// 	2	ASK,
	// 	3	ASK_SIZE,
	// 	4	DAILY_CHANGE,
	// 	5	DAILY_CHANGE_PERC,
	// 	6	LAST_PRICE,
	// 	7	VOLUME,
	// 	8	HIGH,
	// 	9	LOW
	// 	]
	// ]
	return &market.Ticker{
		High:      pt.Ticker[8],
		Low:       pt.Ticker[9],
		Vol:       pt.Ticker[7],
		Last:      pt.Ticker[6],
		Buy:       pt.Ticker[2],
		Sell:      pt.Ticker[0],
		Timestamp: uint64(time.Now().UTC().UnixNano()),
		// FIXME CurrencyPair: pt.SymbolPair,
		Market: Name,
	}
}

func (c *TickerConsumer) worker(pairs []SymbolPair) {
	var (
		stream   = c.Consumer.Consume()
		iterator = NewIterator(stream, c.log)

		pairTicker *pairTicker
		err        error
	)

	err = c.preamble(pairs, iterator)
	if err != nil {
		c.log.Error(err)
		return
	}

workerLoop:
	for {
		select {
		case <-c.done:
			break workerLoop
		default:
			pairTicker, err = c.consume(iterator)
			if err != nil {
				switch err.(type) {
				case *ErrContinue:
					continue workerLoop
				default:
					c.log.Error(err)
					return
				}
			}

			c.tickers <- c.convertTicker(pairTicker)
		}
	}
}

func (c *TickerConsumer) Consume(pairs []currencies.CurrencyPair) <-chan *market.Ticker {
	var (
		symbolPairs []SymbolPair
		err         error
	)

	symbolPairs, err = CurrencyPairsToSymbolPairs(c.currencies, pairs)
	if err != nil {
		// FIXME: Thats shit, need fix
		panic(err)
	}

	go c.worker(symbolPairs)

	return c.tickers
}

func (c *TickerConsumer) Close() error {
	close(c.done)
	// XXX: Not closing it, it will be GC'ed
	// Or we could make worker a panic in case of race
	// close(c.tickers)
	return c.Consumer.Close()
}

// FIXME: This is shit, consumer should receive reader by semantic,
// but it can't ATM because consumer subscribes to channels only
// when Consume(...) is called.
func (m *Bitfinex) NewTickerConsumer(c io.ReadWriter) market.TickerConsumer {
	var (
		l = prefixwrapper.New(
			"TickerConsumer: ",
			m.log,
		)
	)

	return &TickerConsumer{
		Consumer: consumer.New(
			wsutil.NewReader(
				c,
				ws.StateClientSide,
			),
			l,
		),
		connection:          c,
		currencies:          m.currencies,
		channelToSymbolPair: symbolPairByChannel{},
		tickers:             make(chan *market.Ticker, 128),
		done:                make(chan struct{}),
		log:                 l,
	}
}
