// Package clic provides a structured multiplexer for CLI commands. In other
// words, clic will parse a CLI command and route callers to the appropriate
// handler.
package clic

import (
	"context"

	"github.com/daved/flagset"
)

// Handler describes types that can be used to handle CLI command requests. Due
// to the nature of CLI commands containing both arguments and flags, a handler
// must expose both a FlagSet along with a HandleCommand function.
type Handler interface {
	FlagSet() *flagset.FlagSet
	HandleCommand(context.Context, *Clic) error
}

// Clic contains a CLI command handler and subcommand handlers.
type Clic struct {
	h        Handler
	subs     []*Clic
	isCalled bool
	parent   *Clic
	meta     map[string]any
}

// New returns a pointer to a newly constructed instance of a Clic.
func New(h Handler, subs ...*Clic) *Clic {
	c := &Clic{
		h:    h,
		subs: subs,
		meta: map[string]any{
			"SkipUsage":   false,
			"SubRequired": false,
		},
	}

	for _, sub := range c.subs {
		sub.parent = c
	}

	return c
}

// Parse receives command line interface arguments. Parse should be run before
// HandleCalled is run or else *Clic cannot know which handler the user
// requires. Parse is a separate function from HandleCalled so that calling code
// can express behavior in between parsing and handling.
func (c *Clic) Parse(args []string) error {
	return parse(c, args, "")
}

// HandleCalled will run the handler that was selected during Parse processing.
func (c *Clic) HandleCalled(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	called := lastCalled(c)
	return called.h.HandleCommand(ctx, called)
}

// Parent returns the parent *Clic if there is one. Otherwise, the value will
// be nil.
func (c *Clic) Parent() *Clic {
	return c.parent
}

// FlagSet exposes the underlying FlagSet being used.
func (c *Clic) FlagSet() *flagset.FlagSet {
	return c.h.FlagSet()
}

// Name returns the name of the current *Clic.
func (c *Clic) Name() string {
	return c.h.FlagSet().Name()
}

// Subs returns all of the current *Clic instance's subcommands.
func (c *Clic) Subs() []*Clic {
	return c.subs
}

// IsCalled returns whether the *Clic was called as part of the user CLI args
// processed by Parse.
func (c *Clic) IsCalled() bool {
	return c.isCalled
}

// Meta returns a map containing values that can be used to control how usage
// output is displayed.
func (c *Clic) Meta() map[string]any {
	return c.meta
}

func parse(c *Clic, args []string, cmd string) error {
	// TODO: validate sub commands, if any
	fs := c.h.FlagSet()

	c.isCalled = cmd == "" || cmd == fs.Name()
	if !c.isCalled {
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

	for _, sub := range c.subs {
		if err := parse(sub, args, cmd); err != nil {
			return err
		}
	}

	return nil
}

func lastCalled(c *Clic) *Clic {
	for _, sub := range c.subs {
		if !sub.isCalled {
			continue
		}

		return lastCalled(sub)
	}

	return c
}
