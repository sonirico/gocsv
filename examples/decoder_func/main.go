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

	stream := bytes.NewReader([]byte("2,\"patito feo\"\n3,\"patito guapo\""))

	cols := []gocsv.ColFactory[duck]{
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
	}

	decoder := gocsv.NewFuncDecoder[duck](
		stream,
		gocsv.NewDecoderOpts(
			gocsv.WithDecoderSeparator(gocsv.SeparatorComma),
		),
		cols...
	)

	err := decoder.Decode(context.Background(), func(_ context.Context, d duck) error {
		log.Println(d.Name, d.Wings)
		return nil
	})

	if err != nil {
		panic(err)
	}
}
