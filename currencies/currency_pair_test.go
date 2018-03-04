package currencies

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurrencyPairJSON(t *testing.T) {
	var (
		samples = []struct {
			name string
			pair CurrencyPair
			json string
		}{
			{
				name: "full",
				pair: CurrencyPair{
					Left:  Currency{Name: "bitcoin", Symbol: "BTC"},
					Right: Currency{Name: "eos", Symbol: "EOS"},
				},
				json: `[{"name":"bitcoin","symbol":"BTC"},{"name":"eos","symbol":"EOS"}]`,
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				buf, err := json.Marshal(sample.pair)
				if err != nil {
					t.Error(err)
					return
				}

				assert.Equal(
					t,
					sample.json,
					string(buf),
				)

				v := CurrencyPair{}
				err = json.Unmarshal(buf, &v)
				if err != nil {
					t.Error(err)
					return
				}

				assert.Equal(
					t,
					sample.pair,
					v,
				)
			},
		)
	}
}
