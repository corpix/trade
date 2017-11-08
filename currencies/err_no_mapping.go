package currencies

import (
	"fmt"
)

type ErrNoMapping struct {
	Currency
	MapperDirection
}

func (e *ErrNoMapping) Error() string {
	return fmt.Sprintf(
		"No currency mapping for '%s(%s)' in direction '%s'",
		e.Currency.Name,
		e.Currency.Symbol,
		e.MapperDirection,
	)
}

func NewErrNoMapping(currency Currency, direction MapperDirection) error {
	return &ErrNoMapping{
		Currency:        currency,
		MapperDirection: direction,
	}
}
