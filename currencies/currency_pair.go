package currencies

import (
	"fmt"
	"strings"

	e "github.com/cryptounicorns/trade/errors"
)

const (
	CurrencyPairDelimiter = "-"
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

func (c CurrencyPair) Eq(cp CurrencyPair) bool {
	return c.Left == cp.Left && c.Right == cp.Right
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

	// FIXME: If mapping will be a map[string]Currency
	// Then we could optimize
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

func CurrencyPairToString(c CurrencyPair, mapping map[Currency]string, delimiter string) (string, error) {
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

	return left + delimiter + right, nil
}

//

func NewCurrencyPair(left, right Currency) CurrencyPair {
	return CurrencyPair{left, right}
}
