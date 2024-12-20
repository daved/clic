package clic

import (
	"fmt"
)

type Error struct {
	cause error
	c     *Clic
}

func NewError(err error, c *Clic) *Error {
	return &Error{err, c}
}

func (e *Error) Error() string {
	return fmt.Sprintf("cli command: %v", e.cause)
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Clic() *Clic {
	return e.c
}
