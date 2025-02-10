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
	"github.com/daved/clic/flagset"
	"github.com/daved/clic/operandset"
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

	handler    Handler
	called     bool
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
		handler:    h,
		Aliases:    names[1:],
		FlagSet:    flagset.New(name),
		OperandSet: operandset.New(name),
		Meta:       make(map[string]any),
	}

	c.Links = Links{
		self: c,
		subs: subs,
	}

	for _, sub := range c.Links.subs {
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
func (c *Clic) Parse(args []string) error {
	applyRecursiveFlags(c.subs, c.FlagSet)

	if err := parseCmdsAndFlags(c, args, c.FlagSet.FlagSet.Name()); err != nil {
		return err
	}

	last := lastCalled(c)
	if err := last.OperandSet.Parse(last.FlagSet.FlagSet.Operands()); err != nil {
		return cerrs.NewError(cerrs.NewParseError(err))
	}

	return nil
}

// HandleResolvedCmd runs the Handler of the command that was selected during
// Parse processing.
func (c *Clic) HandleResolvedCmd(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return c.ResolvedCmd().handler.HandleCommand(ctx)
}

// Flag adds a flag option to the FlagSet. See [flagset.FlagSet.Flag] for
// details about which value types are supported.
func (c *Clic) Flag(val any, names, usage string) *flagset.Flag {
	return c.FlagSet.Flag(val, names, usage)
}

// FlagRecursive adds a flag option to the FlagSet. Recursive flags are not
// visible in the FlagSet instances of child Clic instances before Parse is
// called (i.e. recursive flags are applied to child flagsets in Parse). See
// [flagset.FlagSet.Flag] for details about which value types are supported.
func (c *Clic) FlagRecursive(val any, names, usage string) *flagset.Flag {
	return c.FlagSet.FlagRecursive(val, names, usage)
}

// Operand adds an Operand option to the OperandSet. See
// [operandset.OperandSet.Operand] for details about which value types are
// supported.
func (c *Clic) Operand(val any, req bool, name, desc string) *operandset.Operand {
	return c.OperandSet.Operand(val, req, name, desc)
}

// Usage returns the executed usage template. The Meta fields of the relevant
// types can be leveraged to convey detailed info/behavior in a custom template.
func (c *Clic) Usage() string {
	return NewUsageTmpl(c).String()
}

var errParseNoMatch = errors.New("parse: no command match")

func parseCmdsAndFlags(c *Clic, args []string, cmdName string) (err error) {
	wrap := cerrs.NewError
	fs := c.FlagSet.FlagSet

	c.called = cmdName == "" || cmdName == fs.Name() || slices.Contains(c.Aliases, cmdName)
	if !c.called {
		return errParseNoMatch
	}

	if err := fs.Parse(args); err != nil {
		return wrap(cerrs.NewParseError(err))
	}
	subCmdArgs := fs.Operands()

	if len(subCmdArgs) == 0 {
		if c.SubRequired {
			return wrap(cerrs.NewParseError(ErrSubCmdRequired))
		}

		return nil
	}

	subCmdName := subCmdArgs[0]
	subCmdArgs = subCmdArgs[1:]

	for _, sub := range c.Links.subs {
		if err := parseCmdsAndFlags(sub, subCmdArgs, subCmdName); err != nil {
			if errors.Is(err, errParseNoMatch) {
				continue
			}
			return err
		}
		return nil
	}

	if c.SubRequired {
		return wrap(cerrs.NewParseError(ErrSubCmdRequired))
	}

	return nil
}

func lastCalled(c *Clic) *Clic {
	for _, sub := range c.Links.subs {
		if sub.called {
			return lastCalled(sub)
		}
	}

	return c
}

func resolvedCmdSet(c *Clic) []*Clic {
	all := []*Clic{c}

	for c.parent != nil {
		c = c.parent
		all = append(all, c)
	}

	slices.Reverse(all)

	return all
}

func applyRecursiveFlags(subs []*Clic, src *flagset.FlagSet) {
	for _, sub := range subs {
		flagset.ApplyRecursiveFlags(sub.FlagSet, src)
		applyRecursiveFlags(sub.Links.subs, src)
		applyRecursiveFlags(sub.Links.subs, sub.FlagSet)
	}
}
