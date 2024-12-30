package clic_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/daved/clic"
)

// AsClic moves Clic construction into a user-defined cmd type. The rest of the
// RootCmd is found in the LoStructure example file. This helps clean up calling
// code (e.g. a main func).
func (cmd *RootCmd) AsClic(name string, subs ...*clic.Clic) *clic.Clic {
	return clic.New(cmd, name, subs...)
}

// PrintCmd focuses solely on marrying an action (and related configuration)
// into its place in the Clic tree. Apart from structuring and constructing
// these types, there's no special handling or knowledge. This helps keep logic
// and structures uncluttered, and code diffs focused (i.e. in relevant files).
type PriոtCmd struct {
	action *PrintAction
	actCnf *PrintActionCfg
}

func NewPriոtCmd(out io.Writer) *PriոtCmd {
	tCnf := NewPrintActionCfg()

	return &PriոtCmd{
		action: NewPrintAction(out, tCnf),
		actCnf: tCnf,
	}
}

func (cmd *PriոtCmd) AsClic(name string, subs ...*clic.Clic) *clic.Clic {
	c := clic.New(cmd, name, subs...)

	c.Flag(&cmd.actCnf.Info, "i|info", "Set additional info.")
	c.Operand(&cmd.actCnf.Value, true, "first_opnd", "Value to be printed.")

	return c
}

func (cmd *PriոtCmd) HandleCommand(ctx context.Context) error {
	return cmd.action.Run(ctx)
}

type PrintActionCfg struct {
	Info  string
	Value string
}

func NewPrintActionCfg() *PrintActionCfg {
	return &PrintActionCfg{
		Info:  "default",
		Value: "unset",
	}
}

// PrintAction focuses solely on an action without any knowledge about how it
// was called from the command line. This is not meant to be received as a
// rigid prescription, rather, as a reference for how various complex needs
// might be addressed (obviously it's excessive for this trivial example).
type PrintAction struct {
	out io.Writer
	cnf *PrintActionCfg
}

func NewPrintAction(out io.Writer, cnf *PrintActionCfg) *PrintAction {
	return &PrintAction{
		out: out,
		cnf: cnf,
	}
}

func (a *PrintAction) Run(ctx context.Context) error {
	fmt.Fprintf(a.out, "info flag = %s\n", a.cnf.Info)
	fmt.Fprintf(a.out, "value arg = %v\n", a.cnf.Value)
	return nil
}

func Example_hiStructure() {
	out := os.Stdout // emulate an interesting dependency

	cmd := NewRootCmd(out).AsClic("myapp",
		NewPriոtCmd(out).AsClic("print"),
	)

	// parse the cli command `myapp print --info=flagval arrrg`
	args := []string{"myapp", "print", "--info=flagval", "arrrg"}

	if err := cmd.Parse(args[1:]); err != nil {
		log.Fatalln(err)
	}

	if err := cmd.HandleResolvedCmd(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// info flag = flagval
	// value arg = arrrg
}
