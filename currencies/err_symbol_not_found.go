package currencies

import (
	"fmt"
)

type ErrSymbolNotFound struct {
	Symbol
}

func (e *ErrSymbolNotFound) Error() string {
	return fmt.Sprintf(
		"Currency with symbol '%s' is not found",
		e.Symbol,
	)
}

func NewErrSymbolNotFound(symbol Symbol) error {
	return &ErrSymbolNotFound{
		Symbol: symbol,
	}
}
