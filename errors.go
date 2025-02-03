package clic

import (
	"errors"
	"fmt"

	"github.com/daved/flagset"
	"github.com/daved/operandset"
)

// ErrSubCmdRequired signals that a subcommand is required and not set.
var ErrSubCmdRequired = errors.New("subcommand is required and not set")

// Cause values are provided for documentation, and to allow callers to easily
// detect error conditions using a switch/case and [errors.Is]. If error
// inspection is required, use [errors.As].
var (
	CauseSubCmdRequired   = ErrSubCmdRequired
	CauseFlagHydrate      = &flagset.FlagHydrateError{}
	CauseFlagUnrecognized = &flagset.FlagUnrecognizedError{}
	CauseOperandHydrate   = &operandset.OperandHydrateError{}
	CauseOperandMissing   = &operandset.OperandMissingError{}
)

func UserFriendlyError(err error) error {
	if errors.Is(err, CauseSubCmdRequired) {
		return errors.New("A subcommand is required")
	}
	if fhErr := (*flagset.FlagHydrateError)(nil); errors.As(err, &fhErr) {
		return fmt.Errorf("Bad flag value for %q (%v)", fhErr.Name, fhErr.Unwrap())
	}
	if fuErr := (*flagset.FlagUnrecognizedError)(nil); errors.As(err, &fuErr) {
		return fmt.Errorf("Unrecognized flag %q", fuErr.Name)
	}
	if ohErr := (*operandset.OperandHydrateError)(nil); errors.As(err, &ohErr) {
		return fmt.Errorf("Bad operand value for %q (%v)", ohErr.Name, ohErr.Unwrap())
	}
	if omErr := (*operandset.OperandMissingError)(nil); errors.As(err, &omErr) {
		return fmt.Errorf("Operand %q is required", omErr.Name)
	}
	return err
}
