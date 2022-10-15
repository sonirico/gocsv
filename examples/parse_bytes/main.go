package main

import (
	"github.com/sonirico/gocsv"
	"log"
)

func main() {

	type duck struct {
		Wings int
		Name  string
	}

	payload := []byte(`2,"patito feo"`)

	parser := gocsv.NewParser[duck](
		gocsv.SeparatorComma,
		gocsv.IntCol[duck](
			gocsv.QuoteNone,
			nil,
			func(d *duck, x int) {
				d.Wings = x
			},
		),
		gocsv.StringCol[duck](
			gocsv.QuoteDouble,
			nil,
			func(d *duck, x string) {
				d.Name = x
			},
		),
	)

	patitoFeo := duck{}

	if err := parser.Parse(payload, &patitoFeo); err != nil {
		panic(err)
	}

	log.Println(patitoFeo)
}
