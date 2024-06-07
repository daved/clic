package clic

import (
	"context"

	"github.com/daved/flagset"
)

type Handler interface {
	FlagSet() *flagset.FlagSet
	HandleCommand(context.Context, *Clic) error
}

type Clic struct {
	h        Handler
	subs     []*Clic
	isCalled bool
	parent   *Clic
	meta     map[string]any
}

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

func (c *Clic) Parse(args []string) error {
	return parse(c, args, "")
}

func (c *Clic) HandleCalled(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	called := lastCalled(c)
	return called.h.HandleCommand(ctx, called)
}

func (c *Clic) Parent() *Clic {
	return c.parent
}

func (c *Clic) FlagSet() *flagset.FlagSet {
	return c.h.FlagSet()
}

func (c *Clic) Name() string {
	return c.h.FlagSet().Name()
}

func (c *Clic) Subs() []*Clic {
	return c.subs
}

func (c *Clic) IsCalled() bool {
	return c.isCalled
}

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
