package clic_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/daved/clic"
)

type HandleRoot struct {
	out io.Writer
}

func NewHandleRoot(out io.Writer) *HandleRoot {
	return &HandleRoot{
		out: out,
	}
}

func (s *HandleRoot) HandleCommand(ctx context.Context) error {
	fmt.Fprintln(s.out, "from root")
	return nil
}

type HandlePrint struct {
	out   io.Writer
	info  string
	value string
}

func NewHandlePrint(out io.Writer) *HandlePrint {
	return &HandlePrint{
		out:   out,
		info:  "default",
		value: "unset",
	}
}

func (s *HandlePrint) HandleCommand(ctx context.Context) error {
	fmt.Fprintf(s.out, "info flag = %s\nvalue operand = %v\n", s.info, s.value)
	return nil
}

func Example_loStructure() {
	out := os.Stdout // emulate an interesting dependency

	// Associate Handler with command name
	handlePrint := NewHandlePrint(out)
	print := clic.New(handlePrint, "print")

	// Associate "print" flag and operand variables with relevant names
	print.Flag(&handlePrint.info, "i|info", "Set additional info.")
	print.Operand(&handlePrint.value, true, "first_operand", "Value to be printed.")

	// Associate Handler with application name, adding "print" as a subcommand
	root := clic.New(NewHandleRoot(out), "myapp", print)

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
