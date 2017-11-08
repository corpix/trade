package currencies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMappingFromIntersection(t *testing.T) {
	var (
		samples = []struct {
			name    string
			left    Currencies
			right   Currencies
			mapping Mapping
		}{
			{
				name: "merging direction left less items",
				left: Currencies{
					Currency{"foo", "BAR"},
				},
				right: Currencies{
					Currency{"foo", "BAR"},
					Currency{"baz", "QUX"},
				},
				mapping: Mapping{
					Currency{"foo", "BAR"}: Currency{"foo", "BAR"},
				},
			},
			{
				name: "merging direction right less items",
				left: Currencies{
					Currency{"foo", "BAR"},
					Currency{"baz", "QUX"},
				},
				right: Currencies{
					Currency{"foo", "BAR"},
				},
				mapping: Mapping{
					Currency{"foo", "BAR"}: Currency{"foo", "BAR"},
				},
			},
			{
				name: "different symbols with same name",
				left: Currencies{
					Currency{"foo", "BAR"},
				},
				right: Currencies{
					Currency{"foo", "BAZ"},
				},
				mapping: Mapping{
					Currency{"foo", "BAR"}: Currency{"foo", "BAZ"},
				},
			},
			{
				name: "different symbols with same name reversed",
				left: Currencies{
					Currency{"foo", "BAZ"},
				},
				right: Currencies{
					Currency{"foo", "BAR"},
				},
				mapping: Mapping{
					Currency{"foo", "BAZ"}: Currency{"foo", "BAR"},
				},
			},
		}
	)

	for _, sample := range samples {
		t.Run(
			sample.name,
			func(t *testing.T) {
				assert.Equal(
					t,
					sample.mapping,
					NewMappingFromIntersection(
						sample.left,
						sample.right,
					),
				)
			},
		)
	}
}
