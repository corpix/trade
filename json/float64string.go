package json

import (
	"strconv"
	"strings"
)

type Float64String float64

func (j *Float64String) MarshalJSON() ([]byte, error) {
	return []byte(
		`"` + strconv.FormatFloat(
			float64(*j),
			'f',
			6,
			64,
		) + `"`,
	), nil
}
func (j *Float64String) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseFloat(
		strings.Trim(string(data), `"`),
		64,
	)
	if err != nil {
		return err
	}

	*j = Float64String(v)

	return nil
}
