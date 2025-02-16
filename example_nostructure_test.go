package clic_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/daved/clic"
)

var args = []string{"myapp", "print", "--info=flagval", "arrrg"}

func Example_noStructure() {
	var (
		info  = "default"
		value = "unset"
		out   = os.Stdout // emulate an interesting dependency
	)

	// Associate HandlerFunc with command name "print"
	print := clic.NewFromFunc(
		func(ctx context.Context) error {
			fmt.Fprintf(out, "info flag = %s\nvalue operand = %v\n", info, value)
			return nil
		},
		"print")

	// Associate "print" flag and operand variables with relevant names
	print.Flag(&info, "i|info", "Set additional info.")
	print.Operand(&value, true, "first_operand", "Value to be printed.")

	// Associate HandlerFunc with application name, adding "print" as a subcommand
	root := clic.NewFromFunc(
		func(ctx context.Context) error {
			fmt.Fprintln(out, "ouch, hit root")
			return nil
		},
		"myapp", print)

	// Parse the cli command as `myapp print --info=flagval arrrg`
	resolved, err := root.Parse(args[1:])
	if err != nil {
		log.Fatalln(err)
	}

	// Run the handler that Parse resolved to
	if err := resolved.Handle(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// info flag = flagval
	// value operand = arrrg
}
