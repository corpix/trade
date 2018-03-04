package currencies

import (
	"fmt"
	"strings"
)

type ErrParseSymbolPair struct {
	RawSymbolPair string
	reason        []string
}

func (e *ErrParseSymbolPair) Error() string {
	var (
		r string
	)

	if len(e.reason) > 0 {
		r = fmt.Sprintf(
			", reason '%s'",
			strings.Join(e.reason, ""),
		)
	}

	return fmt.Sprintf(
		"Failed to parse pair '%s'%s",
		e.RawSymbolPair,
		r,
	)
}

func NewErrParseSymbolPair(rawSymbolPair string, reason ...string) error {
	return &ErrParseSymbolPair{
		RawSymbolPair: rawSymbolPair,
		reason:        reason,
	}
}
