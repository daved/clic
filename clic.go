// Package clic provides a structured multiplexer for CLI commands. In other
// words, clic will parse a CLI command and route callers to the appropriate
// handler.
package clic

import (
	"context"

	"github.com/daved/flagset"
)

// MetaKey constants document which keys can be used in a Clic Meta map that
// are leveraged by this package's usage template.
var (
	MetaKeySkipUsage   = "SkipUsage"
	MetaKeySubRequired = "SubRequired"
	MetaKeyCmdDesc     = "CmdDesc"
	MetaKeyArgsHint    = "ArgsHint"
)

// Handler describes types that can be used to handle CLI command requests. Due
// to the nature of CLI commands containing both arguments and flags, a handler
// must expose both a FlagSet along with a HandleCommand function.
type Handler interface {
	FlagSet() *flagset.FlagSet
	HandleCommand(context.Context) error
}

// Clic contains a CLI command handler and subcommand handlers.
type Clic struct {
	Handler  Handler
	Subs     []*Clic
	IsCalled bool
	Parent   *Clic
	Meta     map[string]any
	tmplTxt  string
}

// New returns a pointer to a newly constructed instance of a Clic.
func New(h Handler, subs ...*Clic) *Clic {
	c := &Clic{
		Handler: h,
		Subs:    subs,
		Meta: map[string]any{
			MetaKeySkipUsage:   false,
			MetaKeySubRequired: false,
		},
		tmplTxt: tmplText,
	}

	for _, sub := range c.Subs {
		sub.Parent = c
	}

	return c
}

// SetUsageTemplate allows callers to override the base template text.
func (c *Clic) SetUsageTemplate(txt string) {
	c.tmplTxt = txt
}

// Parse receives command line interface arguments. Parse should be run before
// Called or Handle so that *Clic can know which handler the user requires.
// Parse is a separate function from Called and Handle so that calling code can
// express behavior in between parsing and handling.
func (c *Clic) Parse(args []string) error {
	return parse(c, args, "")
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

func parse(c *Clic, args []string, cmd string) error {
	// TODO: validate sub commands, if any
	fs := c.Handler.FlagSet()

	c.IsCalled = cmd == "" || cmd == fs.Name()
	if !c.IsCalled {
		return nil
	}

	if err := fs.Parse(args); err != nil {
		return NewParseError(err, c)
	}
	args = fs.Args()

	nArg := fs.NArg()
	if nArg == 0 {
		return nil
	}

	cmd = args[len(args)-nArg]
	args = args[len(args)-nArg+1:]

	for _, sub := range c.Subs {
		if err := parse(sub, args, cmd); err != nil {
			return err
		}
	}

	return nil
}

func lastCalled(c *Clic) *Clic {
	for _, sub := range c.Subs {
		if !sub.IsCalled {
			continue
		}

		return lastCalled(sub)
	}

	return c
}
