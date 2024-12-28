package vtypes

import "encoding"

type OperandFunc func(string) error

type TextMarshalUnmarshaler interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}
