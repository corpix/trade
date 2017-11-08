package currencies

import (
	"fmt"
)

type ErrNotFound struct {
	Name string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf(
		"Currency with name '%s' is not found",
		e.Name,
	)
}

func NewErrNotFound(name string) error {
	return &ErrNotFound{
		Name: name,
	}
}
