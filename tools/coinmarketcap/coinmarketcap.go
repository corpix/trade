package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/corpix/loggers/logger/logrus"
	logrusLogger "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	Addr = "https://coinmarketcap.com/"

	noVolume = "None"
)

var (
	log     = logrus.New(logrusLogger.New())
	newLine = []byte{'\n'}
)

var (
	Paths = map[string]string{
		"all":       "all/views/all/",
		"exchanges": "exchanges/",
	}

	Commands = []cli.Command{
		cli.Command{
			Name:   "all",
			Usage:  "Get data about all coins",
			Flags:  []cli.Flag{},
			Action: AllAction,
		},
		cli.Command{
			Name:  "exchanges",
			Usage: "Get data about concrete exchange coins",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "exchange",
					Usage: "Exchange name, you could see available " +
						"exchanges here https://coinmarketcap.com/exchanges/volume/24-hour/",
				},
			},
			Action: ExchangesAction,
		},
	}
)

type Currency struct {
	Name   string  `json:"name"`
	Symbol string  `json:"symbol"`
	Volume float64 `json:"volume"`
}

func write(buf []byte) error {
	var (
		err error
	)

	_, err = os.Stdout.Write(buf)
	if err != nil {
		return err
	}

	_, err = os.Stdout.Write(newLine)
	if err != nil {
		return err
	}

	return nil
}

func parseVolume(volume string) (float64, error) {
	if volume == "" || volume == noVolume {
		volume = "0"
	}

	return strconv.ParseFloat(
		strings.TrimSpace(volume),
		64,
	)
}

func AllAction(ctx *cli.Context) error {
	var (
		d   *goquery.Document
		err error
	)

	d, err = goquery.NewDocument(Addr + Paths["all"])
	if err != nil {
		return err
	}

	d.
		Find("#currencies-all tbody tr").
		EachWithBreak(
			func(_ int, s *goquery.Selection) bool {
				// Nothing to skip here, tbody does not contain header.
				var (
					name      = s.Find(".currency-name > .currency-name-container").Text()
					symbol    = s.Find(".currency-name > .currency-symbol").Text()
					volume, _ = s.Find("td > .volume").Attr("data-usd")

					c = Currency{
						Name:   strings.TrimSpace(name),
						Symbol: strings.TrimSpace(symbol),
					}

					buf []byte
				)

				c.Volume, err = parseVolume(volume)
				if err != nil {
					return false
				}

				buf, err = json.Marshal(&c)
				if err != nil {
					return false
				}

				err = write(buf)
				if err != nil {
					return false
				}

				return true
			},
		)

	return err
}

func ExchangesAction(ctx *cli.Context) error {
	var (
		exchange = ctx.String("exchange")
		d        *goquery.Document
		err      error
	)

	if exchange == "" {
		return errors.New("You should provide an exchange name")
	}

	d, err = goquery.NewDocument(Addr + Paths["exchanges"] + exchange)
	if err != nil {
		return err
	}

	d.
		Find("#markets table tbody tr").
		EachWithBreak(
			func(n int, s *goquery.Selection) bool {
				if n == 0 {
					// Skip first tr in tbody,
					// because it is a table header -_-
					return true
				}

				var (
					name   = s.Find("td > a.market-name").Text()
					symbol = strings.SplitN(
						s.Find("td:nth-child(3)").Text(),
						"/",
						2,
					)[0]
					volume, _ = s.Find("td > .volume").Attr("data-usd")

					c = Currency{
						Name:   strings.TrimSpace(name),
						Symbol: strings.TrimSpace(symbol),
					}

					buf []byte
				)

				c.Volume, err = parseVolume(volume)
				if err != nil {
					return false
				}

				buf, err = json.Marshal(&c)
				if err != nil {
					return false
				}

				err = write(buf)
				if err != nil {
					return false
				}

				return true
			},
		)

	return err
}

func main() {
	var (
		app = cli.NewApp()
		err error
	)

	app.Commands = Commands

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
