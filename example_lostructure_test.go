package clic_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/daved/clic"
)

type RootHandle struct {
	out io.Writer
}

func NewRootHandle(out io.Writer) *RootHandle {
	return &RootHandle{
		out: out,
	}
}

func (s *RootHandle) HandleCommand(ctx context.Context) error {
	fmt.Fprintln(s.out, "ouch, hit root")
	return nil
}

type PrintHandle struct {
	out   io.Writer
	info  string
	value string
}

func NewPrintHandle(out io.Writer) *PrintHandle {
	return &PrintHandle{
		out:   out,
		info:  "default",
		value: "unset",
	}
}

func (s *PrintHandle) HandleCommand(ctx context.Context) error {
	fmt.Fprintf(s.out, "info flag = %s\n", s.info)
	fmt.Fprintf(s.out, "value arg = %v\n", s.value)
	return nil
}

func Example_loStructure() {
	out := os.Stdout // emulate an interesting dependency

	// Associate Handler with command name "print"
	printHandle := NewPrintHandle(out)
	print := clic.New(printHandle, "print")

	// Associate "print" flag and operand variables with relevant names
	print.Flag(&printHandle.info, "i|info", "Set additional info.")
	print.Operand(&printHandle.value, true, "first_operand", "Value to be printed.")

	// Associate Handler with application name, adding "print" as a subcommand
	rootHandle := clic.New(NewRootHandle(out), "myapp", print)

	// Parse the cli command as `myapp print --info=flagval arrrg`
	if err := rootHandle.Parse(args[1:]); err != nil {
		log.Fatalln(err)
	}

	// Run the handler that Parse resolved to
	if err := rootHandle.HandleResolvedCmd(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// info flag = flagval
	// value arg = arrrg
}
