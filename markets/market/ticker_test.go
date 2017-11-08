package market

import (
	"testing"

	"github.com/corpix/formats"

	"github.com/stretchr/testify/assert"
)

func TestTickerJSON(t *testing.T) {
	samples := []struct {
		name   string
		ticker *Ticker
		result string
		err    error
	}{
		{
			"empty",
			&Ticker{},
			`{` +
				`"high":0,` +
				`"low":0,` +
				`"vol":0,` +
				`"last":0,` +
				`"buy":0,` +
				`"sell":0,` +
				`"timestamp":0,` +
				`"currencyPair":"-",` +
				`"market":""` +
				`}`,
			nil,
		},
	}

	format := formats.NewJSON()

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				result, err := format.Marshal(sample.ticker)
				assert.EqualValues(t, sample.err, err)
				assert.EqualValues(t, sample.result, string(result))
			},
		)
	}
}
