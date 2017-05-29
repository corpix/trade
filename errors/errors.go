package errors

// The MIT License (MIT)
//
// Copyright © 2017 Dmitry Moskowski
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

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
