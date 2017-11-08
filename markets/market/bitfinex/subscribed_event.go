package bitfinex

const (
	SubscribedEventName = "subscribed"
)

type SubscribedEvent struct {
	Event

	Channel string `json:"channel"`
	ChanID  uint   `json:"chanId"`
	Symbol  string `json:"symbol"`
	Pair    string `json:"pair"`
}
