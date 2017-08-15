package cex

import (
	"fmt"
)

// ErrApi is an error indicating that market `m` API
// returned an error.
type ErrApi struct {
	m string
}

func (e *ErrApi) Error() string {
	return fmt.Sprintf(
		"CEX API error '%s'",
		e.m,
	)
}

// NewErrApi creates new ErrApi.
func NewErrApi(m string) error {
	return &ErrApi{m}
}

//
