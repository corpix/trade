package bitfinex

import (
	"io"
	"strconv"
	"time"

	"github.com/corpix/loggers"
	"github.com/corpix/loggers/logger/prefixwrapper"
	"github.com/cryptounicorns/queues/queue/readwriter"
	"github.com/cryptounicorns/queues/result"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"github.com/cryptounicorns/trade/currencies"
	"github.com/cryptounicorns/trade/markets/ticker"
)

const (
	TickerChannelName = "ticker"
)

type TickerConsumer struct {
	*readwriter.Consumer

	connection          io.ReadWriter
	currencies          currencies.Mapper
	channelToSymbolPair symbolPairByChannel
	bufSize             uint
	stream              chan ticker.Result
	done                chan struct{}
	log                 loggers.Logger
}

func (c *TickerConsumer) subscribe(pair SymbolPair, iterator *Iterator) (uint, error) {
	var (
		event          = Event{Event: SubscribeEventName}
		subscribeEvent = SubscribeEvent{
			Event:   event,
			Channel: TickerChannelName,
		}
		subscribeTickerEvent = SubscribeTickerEvent{
			SubscribeEvent: subscribeEvent,
			Pair:           pair,
		}
		e   []byte
		err error
	)

	e, err = Format.Marshal(&subscribeTickerEvent)
	if err != nil {
		return 0, err
	}

	c.log.Debugf(
		"Subscribing to '%+v' on '%s'",
		pair,
		Name,
	)

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
			subscribeTickerEvent,
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
		c.log.Debugf("Subscribed to '%s' with channel id '%+v'", pair, channelID)
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

func (c *TickerConsumer) convertTicker(pt *pairTicker) (*ticker.Ticker, error) {
	var (
		pair currencies.SymbolPair
		err  error
	)

	pair, err = SymbolPairToCommonSymbolPair(
		c.currencies,
		pt.SymbolPair,
	)
	if err != nil {
		return nil, err
	}

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
	return &ticker.Ticker{
		High:       pt.Ticker[8],
		Low:        pt.Ticker[9],
		Vol:        pt.Ticker[7],
		Last:       pt.Ticker[6],
		Buy:        pt.Ticker[2],
		Sell:       pt.Ticker[0],
		Timestamp:  uint64(time.Now().UTC().UnixNano()),
		SymbolPair: pair,
		Market:     Name,
	}, nil
}

func (c *TickerConsumer) worker(pairs []SymbolPair) {
	var (
		stream   <-chan result.Result
		iterator *Iterator

		t          *ticker.Ticker
		pairTicker *pairTicker
		err        error
	)

	stream, err = c.Consumer.Consume()
	if err != nil {
		c.stream <- ticker.Result{Err: err}
		return
	}

	iterator = NewIterator(stream, c.log)

	err = c.preamble(pairs, iterator)
	if err != nil {
		c.stream <- ticker.Result{Err: err}
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
					c.stream <- ticker.Result{Err: err}
					return
				}
			}

			t, err = c.convertTicker(pairTicker)
			if err != nil {
				c.stream <- ticker.Result{Err: err}
				return
			}

			c.stream <- ticker.Result{Value: t}
		}
	}
}

func (c *TickerConsumer) Consume(pairs []currencies.CurrencyPair) (<-chan ticker.Result, error) {
	var (
		symbolPairs []SymbolPair
		err         error
	)

	symbolPairs, err = CurrencyPairsToMarketSymbolPairs(c.currencies, pairs)
	if err != nil {
		return nil, err
	}

	go c.worker(symbolPairs)

	return c.stream, nil
}

func (c *TickerConsumer) Close() error {
	close(c.done)
	// XXX: Not closing it, it will be GC'ed
	// Or we could make worker a panic in case of race
	// close(c.stream)
	return c.Consumer.Close()
}

// XXX: This is shit, consumer should receive reader by semantic,
// but it can't ATM because consumer subscribes to channels only
// when Consume(...) is called.
func (m *Bitfinex) NewTickerConsumer(connection io.ReadWriter) (ticker.Consumer, error) {
	var (
		bufSize = uint(128)
		l       = prefixwrapper.New(
			"TickerConsumer: ",
			m.log,
		)
		consumer *readwriter.Consumer
		err      error
	)

	consumer, err = readwriter.NewConsumer(
		wsutil.NewReader(
			connection,
			ws.StateClientSide,
		),
		readwriter.Config{ConsumerBufferSize: bufSize},
		l,
	)
	if err != nil {
		return nil, err
	}

	return &TickerConsumer{
		Consumer:            consumer,
		connection:          connection,
		currencies:          m.currencies,
		channelToSymbolPair: symbolPairByChannel{},
		bufSize:             bufSize,
		stream:              make(chan ticker.Result, bufSize),
		done:                make(chan struct{}),
		log:                 l,
	}, nil
}
