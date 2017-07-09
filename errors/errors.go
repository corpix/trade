package errors

import (
	"fmt"
)

// ErrArgumentIsNil is an error indicating that invalid values was passed.
type ErrArgumentIsNil struct {
	v interface{}
}

func (e *ErrArgumentIsNil) Error() string {
	return fmt.Sprintf(
		"Argument '%T' is nil",
		e.v,
	)
}

// NewErrArgumentIsNil creates new ErrArgumentIsNil.
func NewErrArgumentIsNil(v interface{}) error {
	return &ErrArgumentIsNil{v}
}

//

// ErrNoCurrencyRepresentation is an error indicating that currency
// has no representation to be used in concrete api calls.
type ErrNoCurrencyRepresentation struct {
	currency string
}

func (e *ErrNoCurrencyRepresentation) Error() string {
	return fmt.Sprintf(
		"No currency representation for '%s'",
		e.currency,
	)
}

// NewErrNoCurrencyRepresentation creates new ErrNoCurrencyRepresentation.
func NewErrNoCurrencyRepresentation(currency string) error {
	return &ErrNoCurrencyRepresentation{
		currency,
	}
}

//

// ErrEndpoint is an error indicating that request to endpoint
// resulted in error.
type ErrEndpoint struct {
	url   string
	error string
	code  int
	want  int
}

func (e *ErrEndpoint) Error() string {
	return fmt.Sprintf(
		"Endpoint request to '%s' finished with error '%s' and code '%d' instead of '%d'",
		e.url,
		e.error,
		e.code,
		e.want,
	)
}

// NewErrEndpoint creates new ErrEndpoint.
func NewErrEndpoint(url string, error string, code int, want int) error {
	return &ErrEndpoint{
		url,
		error,
		code,
		want,
	}
}

//
