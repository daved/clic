package clic

import (
	"fmt"

	"github.com/daved/clic/cerrs"
	"github.com/daved/flagset"
	"github.com/daved/operandset"
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

type (
	ParseError       = cerrs.ParseError
	SubRequiredError = cerrs.SubRequiredError
)

var (
	CauseParse       = &ParseError{}
	CauseSubRequired = &SubRequiredError{}

	CauseFlagSet      = &flagset.Error{}
	CauseParseFlagSet = &flagset.ParseError{}

	CauseOperandSet        = &operandset.Error{}
	CauseParseOperand      = &operandset.ParseError{}
	CauseOperandMissing    = &operandset.OperandMissingError{}
	CauseConvertRawOperand = &operandset.ConvertRawError{}
)
