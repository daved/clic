package clic

import (
	"fmt"

	"github.com/daved/clic/cerrs"
	"github.com/daved/flagset"
	"github.com/daved/operandset"
)

// Error is the package-level error implementation.
type Error struct {
	cause error
	c     *Clic
}

// NewError returns a new instance of Error.
func NewError(err error, c *Clic) *Error {
	return &Error{err, c}
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("cli command: %v", e.cause)
}

// Unwrap implements the [errors] Unwrap anonymous interface.
func (e *Error) Unwrap() error {
	return e.cause
}

// Clic returns the instance of Clic within the Clic tree which encountered
// an error state. This can be useful for taking detailed action when parsing
// fails.
func (e *Error) Clic() *Clic {
	return e.c
}

// Error types forward basic error types from the cerrs package for access and
// documentation. If an error has interesting behavior, it should be defined
// directly in this package.
type (
	ParseError       = cerrs.ParseError
	SubRequiredError = cerrs.SubRequiredError
)

// Cause types allow callers to easily check for potential error values easily
// using a switch/case and [errors.Is].
var (
	CauseParse       = &ParseError{}
	CauseSubRequired = &SubRequiredError{}

	CauseFlagSet      = &flagset.Error{}
	CauseParseFlagSet = &flagset.ParseError{}

	CauseOperandSet        = &operandset.Error{}
	CauseParseOperandSet   = &operandset.ParseError{}
	CauseOperandMissing    = &operandset.OperandMissingError{}
	CauseConvertRawOperand = &operandset.ConvertRawError{}
)
