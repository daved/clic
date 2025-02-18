// Package clic provides a structured multiplexer for CLI commands. In other
// words, clic will parse CLI command arguments and route callers to the
// appropriate handler.
//
// There are three kinds of command line arguments that clic helps to manage:
// Commands/Subcommands, Flags (plus related flag values), and Operands.
// Commands/subcommands each optionally have their own flags and operands. If
// an argument of a command does not match a subcommand, and is not a flag arg
// (i.e. it does not start with a hyphen and is not a flag value), then it will
// be parsed as an operand if any operands have been defined.
//
// Argument kinds and their placements:
//
//	command --flag=flag-value subcommand -f flag-value operand_a operand_b
//
// Custom templates and template behaviors (i.e. template function maps) can be
// set. Custom data can be attached to instances of Clic, FlagSet, Flag,
// OperandSet, and Operand using their Meta fields for access from custom
// templates.
package clic

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/daved/clic/cerrs"
	"github.com/daved/flagset"
	"github.com/daved/operandset"
)

// Handler describes types that can be used to handle CLI command requests.
type Handler interface {
	HandleCommand(context.Context) error
}

// HandlerFunc can be used to easily convert a compatible function to a Handler
// implementation.
type HandlerFunc func(context.Context) error

// HandleCommand implements Handler.
func (f HandlerFunc) HandleCommand(ctx context.Context) error {
	return f(ctx)
}

// Clic manages a CLI command handler and related information.
type Clic struct {
	Links

	Handler    Handler
	Aliases    []string
	FlagSet    *flagset.FlagSet
	OperandSet *operandset.OperandSet

	SubRequired bool
	Description string
	HideUsage   bool

	Category       string
	SubCmdCatsSort []string

	Meta map[string]any
}

// New returns an instance of Clic.
func New(h Handler, name string, subs ...*Clic) *Clic {
	names := strings.Split(name, "|")
	name = names[0]

	c := &Clic{
		Handler:    h,
		Aliases:    names[1:],
		FlagSet:    flagset.New(name),
		OperandSet: operandset.New(name),
		Meta:       make(map[string]any),
	}

	c.Links = Links{
		subs: subs,
	}

	for _, sub := range c.subs {
		sub.parent = c
	}

	return c
}

// NewFromFunc returns an instance of Clic. Any function that is compatible with
// [HandlerFunc] will be converted automatically and used as the [Handler].
func NewFromFunc(f HandlerFunc, name string, subs ...*Clic) *Clic {
	return New(f, name, subs...)
}

// Parse receives command line interface arguments, and should be run before
// HandleResolvedCmd or Link fields are used. Parse is intended to be its own
// step when using Clic so that calling code can express behavior in between
// parsing and handling.
func (c *Clic) Parse(args []string) (*Clic, error) {
	resolved, err := parseCmdsAndFlags(c, args, c.FlagSet.Name())
	if err != nil {
		return resolved, err
	}

	if err := resolved.OperandSet.Parse(resolved.FlagSet.Operands()); err != nil {
		return resolved, cerrs.NewError(cerrs.NewParseError(err))
	}

	return resolved, nil
}

// Flag adds a flag option to the FlagSet. See [flagset.FlagSet.Flag] for
// details about which value types are supported.
func (c *Clic) Flag(val any, names, usage string) *flagset.Flag {
	return c.FlagSet.Flag(val, names, usage)
}

// Operand adds an Operand option to the OperandSet. See
// [operandset.OperandSet.Operand] for details about which value types are
// supported.
func (c *Clic) Operand(val any, req bool, name, desc string) *operandset.Operand {
	return c.OperandSet.Operand(val, req, name, desc)
}

// Recursively applies the provided function to the current Clic instance, and
// all its subcommands recursively.
func (c *Clic) Recursively(fn func(*Clic)) {
	fn(c)
	for _, sub := range c.subs {
		sub.Recursively(fn)
	}
}

func (c *Clic) Handle(ctx context.Context) error {
	return c.Handler.HandleCommand(ctx)
}

// Usage returns the executed usage template. The Meta fields of the relevant
// types can be leveraged to convey detailed info/behavior in a custom template.
func (c *Clic) Usage() string {
	return NewUsageTmpl(c).String()
}

var errParseNoMatch = errors.New("parse: no command match")

func parseCmdsAndFlags(c *Clic, args []string, cmdName string) (*Clic, error) {
	wrap := cerrs.NewError

	called := cmdName == "" || cmdName == c.FlagSet.Name() || slices.Contains(c.Aliases, cmdName)
	if !called {
		return c, errParseNoMatch
	}

	if err := c.FlagSet.Parse(args); err != nil {
		return c, wrap(cerrs.NewParseError(err))
	}
	subCmdArgs := c.FlagSet.Operands()

	if len(subCmdArgs) == 0 {
		if c.SubRequired {
			return c, wrap(cerrs.NewParseError(ErrSubCmdRequired))
		}

		return c, nil
	}

	subCmdName := subCmdArgs[0]
	subCmdArgs = subCmdArgs[1:]

	for _, sub := range c.Links.subs {
		resolved, err := parseCmdsAndFlags(sub, subCmdArgs, subCmdName)
		if err != nil {
			if errors.Is(err, errParseNoMatch) {
				continue
			}
			return resolved, err
		}
		return resolved, nil
	}

	if c.SubRequired {
		return c, wrap(cerrs.NewParseError(ErrSubCmdRequired))
	}

	return c, nil
}
