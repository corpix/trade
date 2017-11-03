package bitfinex

var (
	errContinue = NewErrContinue()
)

type ErrContinue struct{}

func (e *ErrContinue) Error() string {
	return "Loop continuation signal, not an error"
}

func NewErrContinue() *ErrContinue {
	return &ErrContinue{}
}
