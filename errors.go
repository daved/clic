package clic

import (
	"fmt"
)

type ParseError struct {
	cause error
	c     *Clic
}

func NewParseError(err error, c *Clic) *ParseError {
	return &ParseError{err, c}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("cli command: parse: %v", e.cause)
}

func (e *ParseError) Unwrap() error {
	return e.cause
}

func (e *ParseError) Clic() *Clic {
	return e.c
}
