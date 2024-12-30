package clic_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/daved/clic"
)

type RootCmd struct {
	out io.Writer
}

func NewRootCmd(out io.Writer) *RootCmd {
	return &RootCmd{
		out: out,
	}
}

func (cmd *RootCmd) HandleCommand(ctx context.Context) error {
	fmt.Fprintln(cmd.out, "ouch, hit root")
	return nil
}

type PrintCmd struct {
	out   io.Writer
	info  string
	value string
}

func NewPrintCmd(out io.Writer) *PrintCmd {
	return &PrintCmd{
		out:   out,
		info:  "default",
		value: "unset",
	}
}

func (cmd *PrintCmd) HandleCommand(ctx context.Context) error {
	fmt.Fprintf(cmd.out, "info flag = %s\n", cmd.info)
	fmt.Fprintf(cmd.out, "value arg = %v\n", cmd.value)
	return nil
}

func Example_loStructure() {
	out := os.Stdout // emulate an interesting dependency

	printCmd := NewPrintCmd(out)
	printClic := clic.New(printCmd, "print")

	printClic.Flag(&printCmd.info, "i|info", "Set additional info.")
	printClic.Operand(&printCmd.value, true, "first_opnd", "Value to be printed.")

	rootClic := clic.New(NewRootCmd(out), "myapp", printClic)

	// parse the cli command `myapp print --info=flagval arrrg`
	args := []string{"myapp", "print", "--info=flagval", "arrrg"}

	if err := rootClic.Parse(args[1:]); err != nil {
		log.Fatalln(err)
	}

	if err := rootClic.HandleResolvedCmd(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// info flag = flagval
	// value arg = arrrg
}
