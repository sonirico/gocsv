package gocsv

import (
	"bufio"
	"context"
	"io"
)

type (
	stream[T any] interface {
		Next(ctx context.Context) bool
		Data() T
		Err() error
	}

	streamFactory[T any] interface {
		Get(r io.Reader) stream[T]
	}
)

type (
	csvStream[T any] struct {
		current T

		inner stream[[]byte]

		err error

		parser Parser[T]

		ignoreLines int
		ignoreErr   bool

		rowCount int
	}

	newLineStream struct {
		current []byte

		scanner *bufio.Scanner

		err error
	}
)

func (s *newLineStream) Next(_ context.Context) bool {
	if !s.scanner.Scan() {
		s.err = s.scanner.Err()
		return false
	}
	s.current = s.scanner.Bytes()
	return true
}

func (s *newLineStream) Data() []byte {
	return s.current
}

func (s *newLineStream) Err() error {
	if s.err != nil {
		return s.err
	}

	return s.scanner.Err()
}

func (s *csvStream[T]) next(ctx context.Context, shouldParse bool) bool {
	if !s.inner.Next(ctx) {
		return false
	}

	if !shouldParse {
		return true
	}

	var zeroed T
	line := s.inner.Data()
	s.current = zeroed
	s.err = s.parser.Parse(line, &s.current)

	if s.err == nil || s.ignoreErr {
		s.rowCount++
	}

	return true
}

func (s *csvStream[T]) Next(ctx context.Context) bool {
	for s.ignoreLines > 0 {
		_ = s.next(ctx, false)
		s.ignoreLines--
	}
	return s.next(ctx, true)
}

func (s *csvStream[T]) Data() T {
	return s.current
}

func (s *csvStream[T]) Err() error {
	if s.err != nil {
		return s.err
	}

	return s.inner.Err()
}

func newNewLineStream(r io.Reader) stream[[]byte] {
	return &newLineStream{scanner: bufio.NewScanner(r)}
}

func newCSVStream[T any](r io.Reader, ignoreLines int, parser Parser[T]) stream[T] {
	return &csvStream[T]{
		inner:       newNewLineStream(r),
		parser:      parser,
		ignoreLines: ignoreLines,
		ignoreErr:   false,
	}
}
