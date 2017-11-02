package bitfinex

const (
	ErrorEventName = "error"
)

type ErrorEvent struct {
	Event

	Channel string `json:"channel"`
	Symbol  string `json:"symbol"`
	Msg     string `json:"msg"`
	Code    uint   `json:"code"`
	Pair    string `json:"pair"` // FIXME: Use CurrencyPair
}
