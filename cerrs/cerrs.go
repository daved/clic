package cerrs

import (
	"fmt"
	"reflect"
)

// Error is the package-level error implementation.
type Error struct {
	child error
}

// NewError returns a new instance of Error.
func NewError(child error) *Error {
	return &Error{child}
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("cli command: %v", e.child)
}

// Unwrap implements the [errors] Unwrap anonymous interface.
func (e *Error) Unwrap() error {
	return e.child
}

type ParseError struct {
	child error
}

func NewParseError(child error) *ParseError {
	return &ParseError{child}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse: %v", e.child)
}

func (e *ParseError) Unwrap() error {
	return e.child
}

func (e *ParseError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}
