package currencies

import (
	"bytes"
	"fmt"
)

const (
	DefaultSymbolPairDelimiter = "-"
)

var (
	symbolPairJSONShores = []byte{'"'}
)

type SymbolPair struct {
	Left  Symbol
	Right Symbol
}

func (s SymbolPair) MarshalJSON() ([]byte, error) {
	return []byte(
		fmt.Sprintf(
			`"%s%s%s"`,
			s.Left,
			DefaultSymbolPairDelimiter,
			s.Right,
		),
	), nil
}

func (s *SymbolPair) UnmarshalJSON(buf []byte) error {
	var (
		b [][]byte
	)

	b = bytes.Split(
		bytes.TrimSuffix(
			bytes.TrimPrefix(
				bytes.TrimSpace(buf),
				symbolPairJSONShores,
			),
			symbolPairJSONShores,
		),
		[]byte(DefaultSymbolPairDelimiter),
	)

	if len(b) < 2 || len(b[0]) == 0 || len(b[1]) == 0 {
		return NewErrParseSymbolPair(
			string(buf),
			"too few symbols",
		)
	}
	if len(b) > 2 {
		return NewErrParseSymbolPair(
			string(buf),
			"too much symbols",
		)
	}

	s.Left = Symbol(b[0])
	s.Right = Symbol(b[1])

	return nil
}

func FromCurrencyPair(cp CurrencyPair) SymbolPair {
	return SymbolPair{
		Left:  cp.Left.Symbol,
		Right: cp.Right.Symbol,
	}
}

func NewSymbolPair(left Symbol, right Symbol) SymbolPair {
	return SymbolPair{
		Left:  left,
		Right: right,
	}
}
