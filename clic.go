// Package clic provides a structured multiplexer for CLI commands. In other
// words, clic will parse a CLI command and route callers to the appropriate
// handler.
package clic

import (
	"context"
	"errors"
	"slices"

	"github.com/daved/clic/cerrs"
	"github.com/daved/clic/flagset"
	"github.com/daved/clic/operandset"
)

// Handler describes types that can be used to handle CLI command requests. Due
// to the nature of CLI commands containing both arguments and flags, a handler
// must expose both a FlagSet along with a HandleCommand function.
type Handler interface {
	HandleCommand(context.Context) error
}

type HandlerFunc func(context.Context) error

func (f HandlerFunc) HandleCommand(ctx context.Context) error {
	return f(ctx)
}

type Links struct {
	self   *Clic
	parent *Clic
	subs   []*Clic
}

// Clic contains a CLI command handler and subcommand handlers.
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

// New returns a pointer to a newly constructed instance of a Clic.
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

func NewFromFunc(f HandlerFunc, name string, subs ...*Clic) *Clic {
	return New(f, name, subs...)
}

// Parse receives command line interface arguments. Parse should be run before
// Called or Handle so that *Clic can know which handler the user requires.
// Parse is a separate function from Called and Handle so that calling code can
// express behavior in between parsing and handling.
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

// ResolvedCmd returns the command that was selected during Parse processing.
func (l Links) ResolvedCmd() *Clic {
	return lastCalled(l.self)
}

func (l Links) SubCmds() []*Clic {
	return l.subs
}

func (l Links) ParentCmd() *Clic {
	return l.parent
}

// HandleResolvedCmd runs the Handler of the command that was selected during Parse
// processing.
func (c *Clic) HandleResolvedCmd(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return c.ResolvedCmd().handler.HandleCommand(ctx)
}

func (c *Clic) Usage() string {
	data := &TmplData{
		ResolvedCmd:    c,
		ResolvedCmdSet: resolvedCmdSet(c),
	}
	return executeTmpl(c.tmplCfg, data)
}

func (c *Clic) Flag(val any, names, usage string) *flagset.Flag {
	return c.FlagSet.Flag(val, names, usage)
}

func (c *Clic) FlagRecursive(val any, names, usage string) *flagset.Flag {
	return c.FlagSet.FlagRecursive(val, names, usage)
}

func (c *Clic) Operand(val any, req bool, name, desc string) *operandset.Operand {
	return c.OperandSet.Operand(val, req, name, desc)
}

func (c *Clic) SetUsageTemplating(tmplCfg *TmplConfig) {
	c.tmplCfg = tmplCfg
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
