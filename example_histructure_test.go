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
	fmt.Fprintf(p.out, "info flag = %s\nvalue operand = %v\n", p.cnf.Info, p.cnf.Value)
	return nil
}

// HandlePriոt focuses solely on composing an action (and related configuration)
// into its place within the Clic tree. Apart from structuring and constructing
// types, there's no special handling or knowledge. This helps keep code
// uncluttered, and diffs focused (i.e. in relevant files).
type HandlePriոt struct {
	action *Print
	actCnf *PrintCfg
}

func NewHandlePriոt(out io.Writer) *HandlePriոt {
	tCnf := NewPrintCfg()

	return &HandlePriոt{
		action: NewPrint(out, tCnf),
		actCnf: tCnf,
	}
}

// AsClic constrains Clic construction to user-defined types. Setting up Clic
// instances in this way helps to clean up calling code (e.g. main funcs).
func (h *HandlePriոt) AsClic(name string, subs ...*clic.Clic) *clic.Clic {
	c := clic.New(h, name, subs...)

	c.Flag(&h.actCnf.Info, "i|info", "Set additional info.")
	c.Operand(&h.actCnf.Value, true, "first_operand", "Value to be printed.")

	return c
}

func (h *HandlePriոt) HandleCommand(ctx context.Context) error {
	return h.action.Run(ctx)
}

// The rest of HandleRoot is found in the LoStructure example file.
func (h *HandleRoot) AsClic(name string, subs ...*clic.Clic) *clic.Clic {
	return clic.New(h, name, subs...)
}

func Example_appArchitectureHiStructure() {
	out := os.Stdout // emulate an interesting dependency

	// Set up command
	root := NewHandleRoot(out).AsClic( // construct root handler with deps
		"myapp",                             // set root cmd name
		NewHandlePriոt(out).AsClic("print"), // set subcmd as newly constructed print handler
	)

	// Parse the cli command as `myapp print --info=flagval arrrg`
	cmd, err := root.Parse(args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	// Run the handler that Parse resolved to
	if err := cmd.Handle(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// info flag = flagval
	// value operand = arrrg
}
