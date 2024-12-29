package clic

import (
	"fmt"

	"github.com/daved/clic/cerrs"
	"github.com/daved/flagset"
	"github.com/daved/flagset/fserrs"
	"github.com/daved/operandset"
	"github.com/daved/operandset/oserrs"
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

var (
	CauseParse       = &cerrs.ParseError{}
	CauseSubRequired = &cerrs.SubRequiredError{}

	CauseFlagSet      = &flagset.Error{}
	CauseParseFlagSet = &fserrs.ParseError{}

	CauseOperandSet        = &operandset.Error{}
	CauseParseOperand      = &oserrs.ParseError{}
	CauseOperandMissing    = &oserrs.OperandMissingError{}
	CauseConvertRawOperand = &oserrs.ConvertRawError{}
)
