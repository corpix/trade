package bitfinex

const (
	InfoEventName = "info"
)

type InfoEvent struct {
	Event

	Version float64 `json:"version"`
}
