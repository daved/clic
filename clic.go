// Package clic provides a structured multiplexer for CLI commands. In other
// words, clic will parse a CLI command and route callers to the appropriate
// handler.
package clic

import (
	"context"
	"errors"
	"reflect"

	"github.com/daved/flagset"
)

// MetaKey constants document which keys can be used in a Clic Meta map that
// are leveraged by this package's usage template.
var (
	MetaKeySkipUsage   = "SkipUsage"
	MetaKeySubRequired = "SubRequired"
	MetaKeyCmdDesc     = "CmdDesc"
	MetaKeyArgsHint    = "ArgsHint" // TODO: probably replace with args name/desc
)

// Handler describes types that can be used to handle CLI command requests. Due
// to the nature of CLI commands containing both arguments and flags, a handler
// must expose both a FlagSet along with a HandleCommand function.
type Handler interface {
	FlagSet() *flagset.FlagSet
	HandleCommand(context.Context) error
}

type Arg struct {
	Val  any
	Opt  bool
	Name string
	Desc string
}

type ArgSet struct {
	as []*Arg
}

func NewArgSet() *ArgSet {
	return &ArgSet{}
}

func (as *ArgSet) Parse(args []string) error {
	return nil
}

func (as *ArgSet) Arg(val any, opt bool, name, desc string) *ArgSet {
	as.as = append(as.as, &Arg{val, opt, name, desc})
	return as
}

// ArgSetProvider describes types that share an ArgSet.
type ArgSetProvider interface {
	ArgSet() *ArgSet
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
	if _, err := parse(c, args, ""); err != nil {
		return err
	}

	last := lastCalled(c)
	asp, ok := last.Handler.(ArgSetProvider)
	if !ok {
		return nil
	}
	for i, arg := range asp.ArgSet().as {
		v := reflect.ValueOf(arg.Val)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if !v.CanSet() {
			return errors.New("ouch") // TODO: set this
		}

		v.SetString(last.Handler.FlagSet().Arg(i))
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

func parse(c *Clic, args []string, cmd string) (hasSubCalled bool, err error) {
	// TODO: validate sub commands, if any
	fs := c.Handler.FlagSet()

	c.IsCalled = cmd == "" || cmd == fs.Name()
	if !c.IsCalled {
		return false, nil
	}

	if err := fs.Parse(args); err != nil {
		return false, NewParseError(err, c)
	}
	args = fs.Args()

	nArg := fs.NArg()
	if nArg == 0 {
		return false, nil
	}

	cmd = args[len(args)-nArg]
	args = args[len(args)-nArg+1:]

	for _, sub := range c.Subs {
		isSubCalled, subErr := parse(sub, args, cmd)
		if subErr != nil {
			return hasSubCalled, subErr
		}
		if isSubCalled {
			hasSubCalled = true
		}
	}

	return hasSubCalled, nil
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
