package bitfinex

type ErrEmptyDataPayload struct{}

func (e *ErrEmptyDataPayload) Error() string {
	return "Empty data payload"
}

func NewErrEmptyDataPayload() *ErrEmptyDataPayload {
	return &ErrEmptyDataPayload{}
}
