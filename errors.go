package clic

import (
	"errors"

	"github.com/daved/flagset"
	"github.com/daved/operandset"
)

// ErrSubCmdRequired signals that a subcommand is required and not set.
var ErrSubCmdRequired = errors.New("subcommand is required and not set")

// Cause values allow callers to easily check for potential error conditions
// easily using a switch/case and [errors.Is].
var (
	CauseSubCmdRequired   = ErrSubCmdRequired
	CauseFlagHydrate      = &flagset.HydrateError{}
	CauseFlagUnrecognized = &flagset.UnrecognizedFlagError{}
	CauseOperandHydrate   = &operandset.HydrateError{}
	CauseOperandMissing   = &operandset.OperandMissingError{}
)
