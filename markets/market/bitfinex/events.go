package bitfinex

type Event struct {
	Event string `json:"event"`
	Code  int    `json:"code"`
}

type InfoEvent struct {
	Event

	Version float64 `json:"version"`
}
