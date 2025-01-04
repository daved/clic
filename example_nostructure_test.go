package clic_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/daved/clic"
)

func Example_noStructure() {
	var (
		info  = "default"
		value = "unset"
		out   = os.Stdout // emulate an interesting dependency
	)

	// Associate HandlerFunc with command name "print"
	printClic := clic.NewFromFunc(func(ctx context.Context) error {
		fmt.Fprintf(out, "info flag = %s\n", info)
		fmt.Fprintf(out, "value arg = %v\n", value)
		return nil
	}, "print")

	// Associate "print" flag and operand variables with relevant names
	printClic.Flag(&info, "i|info", "Set additional info.")
	printClic.Operand(&value, true, "first_opnd", "Value to be printed.")

	// Associate HandlerFunc with application name, adding "print" as a subcommand
	rootClic := clic.NewFromFunc(func(ctx context.Context) error {
		fmt.Fprintln(out, "ouch, hit root")
		return nil
	}, "myapp", printClic)

	// Parse the cli command as `myapp print --info=flagval arrrg`
	args := []string{"myapp", "print", "--info=flagval", "arrrg"}
	if err := rootClic.Parse(args[1:]); err != nil {
		log.Fatalln(err)
	}

	// Run the handler that Parse resolved to
	if err := rootClic.HandleResolvedCmd(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// info flag = flagval
	// value arg = arrrg
}
