package gocsv

import (
	"bytes"
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestSliceDecoder_Decode(t *testing.T) {
	type (
		duck struct {
			Wings int
			Name  string
		}

		args struct {
			data string
			opts DecodeOpts
			cols []ColFactory[duck]
		}

		want struct {
			expected []duck
			err      error
		}

		testCase struct {
			name string
			args args
			want want
		}
	)

	tests := []testCase{
		{
			name: "no headers, several rows",
			args: args{
				data: "1,sonirico\n2,f3r",
				opts: DecodeOpts{
					ignoreLines: 0,
					separator:   SeparatorComma,
				},
				cols: []ColFactory[duck]{
					IntCol[duck](
						QuoteNone,
						nil,
						func(d *duck, x int) {
							d.Wings = x
						},
					),
					StringCol[duck](
						QuoteNone,
						nil,
						func(d *duck, x string) {
							d.Name = x
						},
					),
				},
			},
			want: want{
				expected: []duck{
					{
						Name:  "sonirico",
						Wings: 1,
					},
					{
						Name:  "f3r",
						Wings: 2,
					},
				},
				err: nil,
			},
		},
		{
			name: "headers, several rows",
			args: args{
				data: "wings,name\n1,sonirico\n2,f3r",
				opts: DecodeOpts{
					ignoreHeaders: true,
					separator:     SeparatorComma,
				},
				cols: []ColFactory[duck]{
					IntCol[duck](
						QuoteNone,
						nil,
						func(d *duck, x int) {
							d.Wings = x
						},
					),
					StringCol[duck](
						QuoteNone,
						nil,
						func(d *duck, x string) {
							d.Name = x
						},
					),
				},
			},
			want: want{
				expected: []duck{
					{
						Name:  "sonirico",
						Wings: 1,
					},
					{
						Name:  "f3r",
						Wings: 2,
					},
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewReader([]byte(test.args.data))

			decoder := NewSliceDecoder[duck](
				buf,
				test.args.opts,
				test.args.cols...,
			)

			var (
				actual []duck
				err    error
			)

			actual, err = decoder.Decode(context.Background())

			if !errors.Is(err, test.want.err) {
				t.Fatalf("unexpected error, want %s, have %s", test.want.err, err)
			}

			if !reflect.DeepEqual(actual, test.want.expected) {
				t.Fatalf("unexpected result, want %v, have %v", test.want.expected, actual)
			}
		})
	}

}
