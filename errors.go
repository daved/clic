package clic

import (
	"fmt"
)

type ParseError struct {
	cause error
	h     Handler
}

func NewParseError(err error, h Handler) *ParseError {
	return &ParseError{err, h}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("command parse: %v", e.cause)
}

func (e *ParseError) Unwrap() error {
	return e.cause
}

func (e *ParseError) Handler() Handler {
	return e.h
}
