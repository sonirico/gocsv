package gocsv

import (
	"context"
	"io"
)

type (
	FuncDecoder[T any] struct {
		stream stream[T]
	}

	Func[T any] func(context.Context, T) error
)

func (d FuncDecoder[T]) Decode(ctx context.Context, fn Func[T]) (err error) {
	for d.stream.Next(ctx) {
		if err = d.stream.Err(); err != nil {
			return
		}

		if err = fn(ctx, d.stream.Data()); err != nil {
			return
		}
	}
	return
}

func NewFuncDecoder[T any](
	r io.Reader,
	opts DecoderOpts,
	colDefs ...ColFactory[T],
) FuncDecoder[T] {
	return FuncDecoder[T]{
		stream: newCSVStream[T](r, opts.ignoreLines, newParser[T](opts.separator, colDefs...)),
	}
}
