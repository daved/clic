package clic

import (
	"fmt"

	"github.com/daved/flagset"
)

type ParseError struct {
	cause error
	fs    *flagset.FlagSet
}

func NewParseError(err error, fs *flagset.FlagSet) *ParseError {
	return &ParseError{err, fs}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("command parse: %v", e.cause)
}

func (e *ParseError) Unwrap() error {
	return e.cause
}

func (e *ParseError) FlagSet() *flagset.FlagSet {
	return e.fs
}
