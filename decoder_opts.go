package gocsv

type (
	DecoderOpts struct {
		ignoreHeaders bool
		ignoreLines   int
		separator     byte
	}

	Wrapped[T any] func(t *T)
)

func NewDecoderOpts(wraps ...Wrapped[DecoderOpts]) DecoderOpts {
	res := DecoderOpts{
		separator: SeparatorComma, // for nice defaults
	}

	for _, wrap := range wraps {
		wrap(&res)
	}

	return res
}

func WithDecoderSeparator(sep byte) Wrapped[DecoderOpts] {
	return func(decoderOpts *DecoderOpts) {
		decoderOpts.separator = sep
	}
}

func WithDecoderIgnoreLines(amount int) Wrapped[DecoderOpts] {
	return func(decoderOpts *DecoderOpts) {
		decoderOpts.ignoreLines = amount
	}
}

func WithDecoderIgnoreHeaders() Wrapped[DecoderOpts] {
	return func(decoderOpts *DecoderOpts) {
		if decoderOpts.ignoreLines == 0 {
			decoderOpts.ignoreLines = 1
		}
	}
}
