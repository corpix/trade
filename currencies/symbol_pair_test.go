package currencies

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbolPairJSON(t *testing.T) {
	var (
		samples = []struct {
			name         string
			pair         SymbolPair
			json         string
			marshalErr   error
			unmarshalErr error
		}{
			{
				name: "full",
				pair: SymbolPair{
					Left:  Symbol("BTC"),
					Right: Symbol("EOS"),
				},
				json: `"BTC-EOS"`,
			},
			{
				name: "too few",
				pair: SymbolPair{
					Left: Symbol("BTC"),
				},
				json:         `"BTC-"`,
				unmarshalErr: NewErrParseSymbolPair(`"BTC-"`, "too few symbols"),
			},
			{
				name: "too much",
				pair: SymbolPair{
					Left: Symbol("BTC"),
					// XXX: Delimiter as a part of symbol is not valid thing
					// probably it should be restricted... but I don't think
					// it is the duty of marshaler.
					Right: Symbol("LTC-ETC"),
				},
				json:         `"BTC-LTC-ETC"`,
				unmarshalErr: NewErrParseSymbolPair(`"BTC-LTC-ETC"`, "too much symbols"),
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				buf, err := json.Marshal(sample.pair)
				assert.Equal(t, sample.marshalErr, err)
				if sample.marshalErr == nil {
					assert.Equal(
						t,
						sample.json,
						string(buf),
					)
				}

				v := SymbolPair{}

				err = json.Unmarshal(buf, &v)
				assert.Equal(t, sample.unmarshalErr, err)
				if sample.unmarshalErr == nil {
					assert.Equal(
						t,
						sample.pair,
						v,
					)
				}
			},
		)
	}
}
