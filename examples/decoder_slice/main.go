package main

import (
	"bytes"
	"context"
	"github.com/sonirico/gocsv"
	"log"
)

func main() {

	type duck struct {
		Wings int
		Name  string
	}

	stream := bytes.NewReader([]byte(`2,"patito feo"`))

	decoder := gocsv.NewSliceDecoder[duck](
		stream,
		gocsv.NewDecoderOpts(
			gocsv.WithDecoderSeparator(gocsv.SeparatorComma),
		),
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

	actual, err := decoder.Decode(context.Background())

	if err != nil {
		panic(err)
	}

	log.Println(actual)
}
