package bitfinex

const (
	SubscribeEventName = "subscribe"
)

type SubscribeEvent struct {
	Event

	Channel string `json:"channel"`
}
