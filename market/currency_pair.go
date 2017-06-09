package market

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
	"strings"

	e "github.com/corpix/trade/errors"
)

type CurrencyPair struct {
	Left  Currency
	Right Currency
}

func (c CurrencyPair) String() string {
	return fmt.Sprintf(
		"%s%s%s",
		c.Left,
		CurrencyPairDelimiter,
		c.Right,
	)
}

func (c CurrencyPair) Format(mapping map[Currency]string, delimiter string) (string, error) {
	var (
		left  string
		right string
		ok    bool
	)

	left, ok = mapping[c.Left]
	if !ok {
		return "", e.NewErrNoCurrencyRepresentation(
			c.Left.String(),
		)
	}
	right, ok = mapping[c.Right]
	if !ok {
		return "", e.NewErrNoCurrencyRepresentation(
			c.Right.String(),
		)
	}

	return fmt.Sprintf(
		"%s%s%s",
		left,
		delimiter,
		right,
	), nil

}

func (c *CurrencyPair) MarshalJSON() ([]byte, error) {
	return []byte(`"` + c.String() + `"`), nil
}

func (c *CurrencyPair) UnmarshalJSON(buf []byte) error {
	cp, err := CurrencyPairFromString(
		strings.Trim(string(buf), `"`),
		CurrencyMapping,
		CurrencyPairDelimiter,
	)
	if err != nil {
		return err
	}

	*c = cp

	return nil
}

func (c *CurrencyPair) MarshalBinary() ([]byte, error) {
	return []byte(c.String()), nil
}

func (c *CurrencyPair) UnmarshalBinary(buf []byte) error {
	cp, err := CurrencyPairFromString(
		string(buf),
		CurrencyMapping,
		CurrencyPairDelimiter,
	)
	if err != nil {
		return err
	}

	*c = cp

	return nil
}

func CurrencyPairFromString(s string, mapping map[Currency]string, delimiter string) (CurrencyPair, error) {
	var (
		left  Currency
		right Currency
	)

	pair := strings.SplitN(s, delimiter, 2)
	if len(pair) != 2 {
		return CurrencyPair{}, e.NewErrNoCurrencyRepresentation(s)
	}
	for k, v := range mapping {
		if pair[0] == v {
			left = k
		}
		if pair[1] == v {
			right = k
		}
	}
	if left == InvalidCurrency {
		return CurrencyPair{}, e.NewErrNoCurrencyRepresentation(pair[0])
	}
	if right == InvalidCurrency {
		return CurrencyPair{}, e.NewErrNoCurrencyRepresentation(pair[1])
	}

	return CurrencyPair{
		left,
		right,
	}, nil
}

//

func NewCurrencyPair(left, right Currency) CurrencyPair {
	return CurrencyPair{left, right}
}
