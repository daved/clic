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

	errs "github.com/daved/clic/clicerrs"
	"github.com/daved/flagset"
)

// Handler describes types that can be used to handle CLI command requests. Due
// to the nature of CLI commands containing both arguments and flags, a handler
// must expose both a FlagSet along with a HandleCommand function.
type Handler interface {
	HandleCommand(context.Context) error
}

type UsageConfig struct {
	Skip     bool
	TmplText string
	CmdDesc  string
}

// Clic contains a CLI command handler and subcommand handlers.
type Clic struct {
	Handler     Handler
	FlagSet     *FlagSet
	ArgSet      *ArgSet
	Subs        []*Clic
	Parent      *Clic
	called      bool
	SubRequired bool
	UsageConfig *UsageConfig
	Meta        map[string]any
}

// New returns a pointer to a newly constructed instance of a Clic.
func New(h Handler, name string, subs ...*Clic) *Clic {
	c := &Clic{
		Handler: h,
		FlagSet: &FlagSet{FlagSet: flagset.New(name)},
		ArgSet:  newArgSet(),
		Subs:    subs,
		UsageConfig: &UsageConfig{
			TmplText: tmplText,
		},
		Meta: make(map[string]any),
	}

	for _, sub := range c.Subs {
		sub.Parent = c
	}

	return c
}

var (
	CauseParseInFlagSet   = &errs.FlagSetError{}
	CauseParseInArgSet    = &errs.ArgSetError{}
	CauseParseArgMissing  = &errs.ArgMissingError{}
	CauseParseSubRequired = &errs.SubRequiredError{}
)

// Parse receives command line interface arguments. Parse should be run before
// Called or Handle so that *Clic can know which handler the user requires.
// Parse is a separate function from Called and Handle so that calling code can
// express behavior in between parsing and handling.
func (c *Clic) Parse(args []string) error {
	applyRecursiveOpts(c, nil)

	if err := parse(c, args, c.FlagSet.FlagSet.Name()); err != nil {
		return err
	}

	last := lastCalled(c)
	if err := last.ArgSet.parse(last.FlagSet.FlagSet.Args()); err != nil {
		return NewError(errs.NewParseError(errs.NewFlagSetError(err)), last)
	}

	return nil
}

// Called returns the command that was selected during Parse processing.
func (c *Clic) Called() *Clic {
	return lastCalled(c)
}

// Handle runs the Handler of the command that was selected during Parse
// processing.
func (c *Clic) Handle(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return c.Called().Handler.HandleCommand(ctx)
}

func (c *Clic) Usage() string {
	data := &tmplData{
		Current: c,
		Called:  allCalled(c),
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

var errParseNoMatch = errors.New("parse: no command match")

func parse(c *Clic, args []string, cmdName string) (err error) {
	fs := c.FlagSet.FlagSet

	c.called = cmdName == "" || cmdName == fs.Name()
	if !c.called {
		return errParseNoMatch
	}

	if err := fs.Parse(args); err != nil {
		return NewError(errs.NewParseError(err), c)
	}
	subCmdArgs := fs.Args()

	if len(subCmdArgs) == 0 {
		if c.SubRequired {
			return NewError(errs.NewParseError(errs.NewSubRequiredError()), c)
		}

		return nil
	}

	subCmdName := subCmdArgs[0]
	subCmdArgs = subCmdArgs[1:]

	for _, sub := range c.Subs {
		if err := parse(sub, subCmdArgs, subCmdName); err != nil {
			if errors.Is(err, errParseNoMatch) {
				continue
			}
			return err
		}
		return nil
	}

	if c.SubRequired {
		return NewError(errs.NewParseError(errs.NewSubRequiredError()), c)
	}

	return nil
}

func lastCalled(c *Clic) *Clic {
	for _, sub := range c.Subs {
		if sub.called {
			return lastCalled(sub)
		}
	}

	return c
}

func allCalled(c *Clic) []*Clic {
	all := []*Clic{c}

	for c.Parent != nil {
		c = c.Parent
		all = append(all, c)
	}

	slices.Reverse(all)

	return all
}
