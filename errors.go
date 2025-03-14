package clic

import (
	"errors"
	"fmt"

	"github.com/daved/flagset"
	"github.com/daved/operandset"
	"github.com/daved/vtypes"
)

// ErrSubCmdRequired signals that a subcommand is required and not set.
var ErrSubCmdRequired = errors.New("subcommand required")

// Cause values are provided for documentation, and to allow callers to easily
// detect error conditions using a switch/case and [errors.Is]. If error
// inspection is required, use [errors.As].
var (
	CauseParseSubCmdRequired   = ErrSubCmdRequired
	CauseParseFlagResolve      = &flagset.ResolveError{}
	CauseParseFlagUnrecognized = flagset.ErrFlagUnrecognized
	CauseParseOperandResolve   = &operandset.ResolveError{}
	CauseParseOperandRequired  = operandset.ErrOperandRequired
	CauseParseHydrateError     = &vtypes.HydrateError{}     // from Flag and Operand Resolve
	CauseParseTypeUnsupported  = vtypes.ErrTypeUnsupported  // from Flag and Operand Resolve
	CauseParseValueUnsupported = vtypes.ErrValueUnsupported // from Flag and Operand Resolve
)

// UserFriendlyError returns a new error containing a plain language message.
func UserFriendlyError(err error) error {
	if errors.Is(err, ErrSubCmdRequired) {
		return errors.New("A subcommand is required")
	}

	if resErr := (*flagset.ResolveError)(nil); errors.As(err, &resErr) {
		if errors.Is(resErr, flagset.ErrFlagUnrecognized) {
			return fmt.Errorf("Unrecognized flag %q", resErr.FlagName)
		}
		if hydErr := friendlyHydrateError(resErr, "flag"); hydErr != nil {
			return hydErr
		}
		return fmt.Errorf("Cannot process flag %q (%v)", resErr.FlagName, resErr.Unwrap())
	}

	if resErr := (*operandset.ResolveError)(nil); errors.As(err, &resErr) {
		if errors.Is(resErr, operandset.ErrOperandRequired) {
			return fmt.Errorf("Operand %q is required", resErr.OperandName)
		}
		if hydErr := friendlyHydrateError(resErr, "operand"); hydErr != nil {
			return hydErr
		}
		return fmt.Errorf("Cannot process operand %q (%v)", resErr.OperandName, resErr.Unwrap())
	}

	return err
}

func friendlyHydrateError(err error, typ string) error {
	if hydErr := (*vtypes.HydrateError)(nil); errors.As(err, &hydErr) {
		if errors.Is(hydErr, vtypes.ErrTypeUnsupported) {
			return fmt.Errorf("Unsupported %s value type '%T'", typ, hydErr.Val)
		}
		return fmt.Errorf(
			"Cannot set %s value of type '%T' (%v)",
			typ, hydErr.Val, hydErr.Unwrap(),
		)
	}
	return nil
}
