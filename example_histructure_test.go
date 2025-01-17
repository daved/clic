package clic_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/daved/clic"
)

type PrintCfg struct {
	Info  string
	Value string
}

func NewPrintCfg() *PrintCfg {
	return &PrintCfg{
		Info:  "default",
		Value: "unset",
	}
}

// Print focuses solely on an action without any knowledge about how it is/was
// called from the command line. This is not meant to be received as a rigid
// prescription, rather, as a reference for how various complex needs might be
// addressed (obviously it's excessive for this trivial example).
type Print struct {
	out io.Writer
	cnf *PrintCfg
}

func NewPrint(out io.Writer, cnf *PrintCfg) *Print {
	return &Print{
		out: out,
		cnf: cnf,
	}
}

func (p *Print) Run(ctx context.Context) error {
	fmt.Fprintf(p.out, "info flag = %s\n", p.cnf.Info)
	fmt.Fprintf(p.out, "value arg = %v\n", p.cnf.Value)
	return nil
}

// PriոtHandle focuses solely on marrying an action (and related configuration)
// into its place within the Clic tree. Apart from structuring and constructing
// these types, there's no special handling or knowledge. This helps keep logic
// and structures uncluttered, and code diffs focused (i.e. in relevant files).
type PriոtHandle struct {
	action *Print
	actCnf *PrintCfg
}

func NewPriոtHandle(out io.Writer) *PriոtHandle {
	tCnf := NewPrintCfg()

	return &PriոtHandle{
		action: NewPrint(out, tCnf),
		actCnf: tCnf,
	}
}

// AsClic moves Clic construction into a user-defined cmd type. Setting up clic
// instances in this way helps to clean up calling code (e.g. main funcs).
func (h *PriոtHandle) AsClic(name string, subs ...*clic.Clic) *clic.Clic {
	c := clic.New(h, name, subs...)

	c.Flag(&h.actCnf.Info, "i|info", "Set additional info.")
	c.Operand(&h.actCnf.Value, true, "first_opnd", "Value to be printed.")

	return c
}

func (h *PriոtHandle) HandleCommand(ctx context.Context) error {
	return h.action.Run(ctx)
}

// The rest of the RootHandle is found in the LoStructure example file.
func (h *RootHandle) AsClic(name string, subs ...*clic.Clic) *clic.Clic {
	return clic.New(h, name, subs...)
}

func Example_hiStructure() {
	out := os.Stdout // emulate an interesting dependency

	// Set up commands (which set up their related actions and clic
	// instances), then relate them using their "AsClic" methods
	c := NewRootHandle(out).AsClic("myapp",
		NewPriոtHandle(out).AsClic("print"),
	)

	// Parse the cli command as `myapp print --info=flagval arrrg`
	if err := c.Parse(args[1:]); err != nil {
		log.Fatalln(err)
	}

	// Run the handler that Parse resolved to
	if err := c.HandleResolvedCmd(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// info flag = flagval
	// value arg = arrrg
}
