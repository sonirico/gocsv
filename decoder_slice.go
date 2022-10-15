package gocsv

import (
	"context"
	"io"
)

type (
	DecodeOpts struct {
		ignoreHeaders bool
		ignoreLines   int
		separator     byte
	}

	SliceDecoder[T any] struct {
		stream stream[T]
	}
)

func (d SliceDecoder[T]) Decode(ctx context.Context) (res []T, err error) {
	for d.stream.Next(ctx) {
		if err = d.stream.Err(); err != nil {
			return
		}

		res = append(res, d.stream.Data())
	}
	return
}

func NewSliceDecoder[T any](
	r io.Reader,
	opts DecodeOpts,
	colDefs ...ColFactory[T],
) SliceDecoder[T] {
	if opts.ignoreLines == 0 && opts.ignoreHeaders {
		opts.ignoreLines = 1
	}

	return SliceDecoder[T]{
		stream: newCSVStream[T](r, opts.ignoreLines, newParser[T](opts.separator, colDefs...)),
	}
}
