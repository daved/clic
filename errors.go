package clic

import (
	"fmt"

	"github.com/daved/clic/clicerrs"
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
	CauseParseInFlagSet   = &clicerrs.FlagSetError{}
	CauseOperandSet       = &operandset.Error{}
	CauseOperandMissing   = &oserrs.OperandMissingError{}
	CauseParseSubRequired = &clicerrs.SubRequiredError{}
)
