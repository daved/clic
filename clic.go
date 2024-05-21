package clic

import (
	"github.com/daved/flagset"
)

type handleFunc func() error

type Handler interface {
	FlagSet() *flagset.FlagSet
	HandleCommand() error
}

type Clic struct {
	h     Handler
	subs  []*Clic
	isSet bool
}

func New(h Handler, subs ...*Clic) *Clic {
	return &Clic{
		h:    h,
		subs: subs,
	}
}

func (c *Clic) Parse(args []string) error {
	return parse(c, args, "")
}

func (c *Clic) HandleCommand() error {
	fn := getFn(c)
	return fn()
}

func parse(c *Clic, args []string, cmd string) error {
	// TODO: validate sub commands, if any
	fs := c.h.FlagSet()

	c.isSet = cmd == "" || cmd == fs.Name()
	if !c.isSet {
		return nil
	}

	if err := fs.Parse(args); err != nil {
		return NewParseError(err, c.h)
	}

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

func getFn(c *Clic) handleFunc {
	for _, sub := range c.subs {
		if !sub.isSet {
			continue
		}

		return getFn(sub)
	}

	return c.h.HandleCommand
}