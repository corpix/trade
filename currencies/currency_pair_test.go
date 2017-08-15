package currencies

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestCurrencyPairFormat(t *testing.T) {
	samples := []struct {
		left   Currency
		right  Currency
		err    error
		output string
	}{
		{
			Bitcoin,
			UnitedStatesDollar,
			nil,
			Bitcoin.String() + CurrencyPairDelimiter + UnitedStatesDollar.String(),
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		v, err := NewCurrencyPair(
			sample.left,
			sample.right,
		).Format(
			CurrencyMapping,
			CurrencyPairDelimiter,
		)
		assert.EqualValues(t, err, sample.err, msg)
		assert.EqualValues(t, v, sample.output, msg)
	}
}

func TestCurrencyPairFromString(t *testing.T) {
	samples := []struct {
		input string
		err   error
		left  Currency
		right Currency
	}{
		{
			Bitcoin.String() + CurrencyPairDelimiter + UnitedStatesDollar.String(),
			nil,
			Bitcoin,
			UnitedStatesDollar,
		},
	}

	for k, sample := range samples {
		msg := spew.Sdump(k, sample)

		v, err := CurrencyPairFromString(
			sample.input,
			CurrencyMapping,
			CurrencyPairDelimiter,
		)
		assert.EqualValues(t, err, sample.err, msg)
		assert.Equal(t, v.Left, sample.left, msg)
		assert.Equal(t, v.Right, sample.right, msg)
	}
}
