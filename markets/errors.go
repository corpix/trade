package markets

import (
	"fmt"
)

// ErrUnsupportedMarket is an error indicating that market `m`
// is not supported in code.
type ErrUnsupportedMarket struct {
	m string
}

func (e *ErrUnsupportedMarket) Error() string {
	return fmt.Sprintf(
		"Unsupported market '%s'",
		e.m,
	)
}

// NewErrUnsupportedMarket creates new ErrUnsupportedMarket.
func NewErrUnsupportedMarket(m string) error {
	return &ErrUnsupportedMarket{m}
}

//
