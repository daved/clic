// Package clic provides a structured multiplexer for CLI commands. In other
// words, clic will parse CLI command arguments and route callers to the
// appropriate handler. Flags and operands can be set up easily.
package clic

import (
	"context"
	"errors"
	"slices"

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
	FlagSet    *flagset.FlagSet
	OperandSet *operandset.OperandSet
	called     bool

	tmplCfg     *TmplConfig
	SubRequired bool
	Description string
	HideUsage   bool
	Meta        map[string]any
}

// New returns an instance of Clic.
func New(h Handler, name string, subs ...*Clic) *Clic {
	c := &Clic{
		handler:    h,
		FlagSet:    flagset.New(name),
		OperandSet: operandset.New(name),
		tmplCfg:    NewDefaultTmplConfig(),
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

	if err := parse(c, args, c.FlagSet.FlagSet.Name()); err != nil {
		return err
	}

	last := lastCalled(c)
	if err := last.OperandSet.Parse(last.FlagSet.FlagSet.Operands()); err != nil {
		return NewError(cerrs.NewParseError(err), last)
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

// SetUsageTemplating is used to override the base template text, and provide a
// custom FuncMap. If a nil FuncMap is provided, no change will be made to the
// existing value.
func (c *Clic) SetUsageTemplating(tmplCfg *TmplConfig) {
	c.tmplCfg = tmplCfg
}

// Usage returns the executed usage template. The Meta fields of the relevant
// types can be leveraged to convey detailed info/behavior in a custom template.
func (c *Clic) Usage() string {
	data := &TmplData{
		ResolvedCmd:    c,
		ResolvedCmdSet: resolvedCmdSet(c),
	}
	return executeTmpl(c.tmplCfg, data)
}

var errParseNoMatch = errors.New("parse: no command match")

func parse(c *Clic, args []string, cmdName string) (err error) {
	fs := c.FlagSet.FlagSet

	c.called = cmdName == "" || cmdName == fs.Name()
	if !c.called {
		return errParseNoMatch
	}

	if err := fs.Parse(args); err != nil {
		return NewError(cerrs.NewParseError(err), c)
	}
	subCmdArgs := fs.Operands()

	if len(subCmdArgs) == 0 {
		if c.SubRequired {
			return NewError(cerrs.NewParseError(cerrs.NewSubRequiredError()), c)
		}

		return nil
	}

	subCmdName := subCmdArgs[0]
	subCmdArgs = subCmdArgs[1:]

	for _, sub := range c.Links.subs {
		if err := parse(sub, subCmdArgs, subCmdName); err != nil {
			if errors.Is(err, errParseNoMatch) {
				continue
			}
			return err
		}
		return nil
	}

	if c.SubRequired {
		return NewError(cerrs.NewParseError(cerrs.NewSubRequiredError()), c)
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
