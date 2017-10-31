package bitfinex

import (
	"github.com/corpix/formats"
)

var (
	Format formats.Format
)

func init() {
	var (
		err error
	)

	Format, err = formats.New("json")
	if err != nil {
		panic(err)
	}
}
