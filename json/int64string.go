package json

import (
	"strconv"
	"strings"
)

type Int64String int64

func (j *Int64String) MarshalJSON() ([]byte, error) {
	return []byte(
		`"` + strconv.FormatInt(
			int64(*j),
			10,
		) + `"`,
	), nil
}
func (j *Int64String) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseInt(
		strings.Trim(string(data), `"`),
		10,
		64,
	)
	if err != nil {
		return err
	}

	*j = Int64String(v)

	return nil
}
