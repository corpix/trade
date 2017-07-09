package market

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
			BTC,
			USD,
			nil,
			"BTC-USD",
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
			"BTC-USD",
			nil,
			BTC,
			USD,
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
