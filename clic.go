package clic

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/daved/clic/cerrs"
	"github.com/daved/flagset"
	"github.com/daved/operandset"
)

// Handler describes types that can be used to handle CLI command requests.
type Handler interface {
	HandleCommand(context.Context) error
}

// HandlerFunc converts compatible functions to a [Handler] implementation.
type HandlerFunc func(context.Context) error

// HandleCommand implements [Handler].
func (f HandlerFunc) HandleCommand(ctx context.Context) error {
	return f(ctx)
}

// Clic manages a CLI command handler and related information.
type Clic struct {
	Links

	// Templating
	Tmpl           *Tmpl // set to NewUsageTmpl by default
	Description    string
	HideUsage      bool
	Category       string
	SubCmdCatsSort []string
	Meta           map[string]any

	// Additional Configuration
	SubRequired bool

	// Reconfiguration (modify as needed)
	Handler Handler
	Aliases []string

	// Accessing (avoid modification)
	FlagSet    *flagset.FlagSet
	OperandSet *operandset.OperandSet
}

// New returns an instance of Clic.
func New(h Handler, name string, subs ...*Clic) *Clic {
	names := strings.Split(name, "|")
	name = names[0]

	c := &Clic{
		Handler:    h,
		Aliases:    names[1:],
		FlagSet:    flagset.New(name),
		OperandSet: operandset.New(name),
		Meta:       make(map[string]any),
	}

	for _, sub := range subs {
		sub.parent = c
	}
	c.Links = Links{subs: subs}

	c.Tmpl = NewUsageTmpl(c)

	return c
}

// NewFromFunc returns an instance of Clic. Any function that is compatible with
// [HandlerFunc] will be converted automatically.
func NewFromFunc(f HandlerFunc, name string, subs ...*Clic) *Clic {
	return New(f, name, subs...)
}

// Parse resolves arguments to the relevant *Clic instance.
func (c *Clic) Parse(args []string) (*Clic, error) {
	resolved, err := parseCmdsAndFlags(c, args, c.FlagSet.Name())
	if err != nil {
		return resolved, err
	}

	if err := resolved.OperandSet.Parse(resolved.FlagSet.Operands()); err != nil {
		return resolved, cerrs.NewError(cerrs.NewParseError(err))
	}

	return resolved, nil
}

// Flag adds a flag option. See [flagset.FlagSet.Flag] for more details like
// compatible value types.
func (c *Clic) Flag(val any, names, usage string) *flagset.Flag {
	return c.FlagSet.Flag(val, names, usage)
}

// Operand adds an operand option. See [operandset.OperandSet.Operand] for more
// details like compatible value types.
func (c *Clic) Operand(val any, req bool, name, desc string) *operandset.Operand {
	return c.OperandSet.Operand(val, req, name, desc)
}

// Recursively applies the provided function to the current Clic instance and
// all its subcommands recursively.
func (c *Clic) Recursively(fn func(*Clic)) {
	fn(c)
	for _, sub := range c.subs {
		sub.Recursively(fn)
	}
}

// Handle calls the HandleCommand method on the set [Handler].
func (c *Clic) Handle(ctx context.Context) error {
	return c.Handler.HandleCommand(ctx)
}

// Usage returns usage text. The default template construction function
// ([NewUsageTmpl]) can be used as a reference for custom templates which should
// be used to set the "Tmpl" field on Clic (likely using [*Clic.Recursively]).
func (c *Clic) Usage() string {
	return c.Tmpl.String()
}

var errParseNoMatch = errors.New("parse: no command match")

func parseCmdsAndFlags(c *Clic, args []string, cmdName string) (*Clic, error) {
	wrap := cerrs.NewError

	called := cmdName == "" || cmdName == c.FlagSet.Name() || slices.Contains(c.Aliases, cmdName)
	if !called {
		return c, errParseNoMatch
	}

	if err := c.FlagSet.Parse(args); err != nil {
		return c, wrap(cerrs.NewParseError(err))
	}
	subCmdArgs := c.FlagSet.Operands()

	if len(subCmdArgs) == 0 {
		if c.SubRequired {
			return c, wrap(cerrs.NewParseError(ErrSubCmdRequired))
		}

		return c, nil
	}

	subCmdName := subCmdArgs[0]
	subCmdArgs = subCmdArgs[1:]

	for _, sub := range c.Links.subs {
		resolved, err := parseCmdsAndFlags(sub, subCmdArgs, subCmdName)
		if err != nil {
			if errors.Is(err, errParseNoMatch) {
				continue
			}
			return resolved, err
		}
		return resolved, nil
	}

	if c.SubRequired {
		return c, wrap(cerrs.NewParseError(ErrSubCmdRequired))
	}

	return c, nil
}
