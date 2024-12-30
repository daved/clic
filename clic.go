// Package clic provides a structured multiplexer for CLI commands. In other
// words, clic will parse a CLI command and route callers to the appropriate
// handler.
package clic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"slices"
	"text/template"

	"github.com/daved/clic/cerrs"
	"github.com/daved/clic/flagset"
	"github.com/daved/operandset"
)

// TODO: consider adding default arg value in tmpl
// TODO: consider adding descriptions of available subcmds to tmpl, optionally skippable
// TODO: consider what subcmd grouping and sorting looks like

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

type UsageConfig struct {
	Skip     bool
	TmplText string
	CmdDesc  string
}

type Links struct {
	self   *Clic
	parent *Clic
	subs   []*Clic
}

// Clic contains a CLI command handler and subcommand handlers.
type Clic struct {
	Links

	// TODO: If it can be updated, expose its field
	// If it needs behavior on setting, expose as setter method
	// TODO: If it should be retrievable for templating and is
	// not exposed as a field, expose getter method
	handler     Handler
	FlagSet     *flagset.FlagSet
	OperandSet  *operandset.OperandSet
	called      bool
	SubRequired bool
	UsageConfig *UsageConfig // TODO: Add Description,HideHint fields / drop UsageConfig
	// TODO: Add templating setter
	Meta map[string]any
	// TODO: Consider renaming Handle(ctx) err to Run or Execute or ?
	// TODO: Consider adding a field of type "Add" which has methods Flag, Operand, FlagRec
	// TODO: Consider renaming "Called"
}

// New returns a pointer to a newly constructed instance of a Clic.
func New(h Handler, name string, subs ...*Clic) *Clic {
	c := &Clic{
		handler:    h,
		FlagSet:    flagset.New(name),
		OperandSet: operandset.New(name),
		UsageConfig: &UsageConfig{
			TmplText: tmplText,
		},
		Meta: make(map[string]any),
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

// CalledCmd returns the command that was selected during Parse processing.
func (l Links) CalledCmd() *Clic {
	return lastCalled(l.self)
}

func (l Links) SubCmds() []*Clic {
	return l.subs
}

func (l Links) ParentCmd() *Clic {
	return l.parent
}

// Handle runs the Handler of the command that was selected during Parse
// processing.
func (c *Clic) Handle(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return c.CalledCmd().handler.HandleCommand(ctx)
}

func (c *Clic) Usage() string {
	data := &tmplData{
		CurrentCmd: c,
		CalledCmds: allCalled(c),
	}

	tmpl := template.New("clic")

	buf := &bytes.Buffer{}

	tmpl, err := tmpl.Parse(c.UsageConfig.TmplText)
	if err != nil {
		fmt.Fprintf(buf, "cli command: template error: %v\n", err)
		return buf.String()
	}

	if err := tmpl.Execute(buf, data); err != nil {
		fmt.Fprintf(buf, "cli command: template error: %v\n", err)
		return buf.String()
	}

	return buf.String()
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

func allCalled(c *Clic) []*Clic {
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
