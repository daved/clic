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

	printCmd := clic.NewFromFunc(func(ctx context.Context) error {
		fmt.Fprintf(out, "info flag = %s\n", info)
		fmt.Fprintf(out, "value arg = %v\n", value)
		return nil
	}, "print")

	printCmd.Flag(&info, "i|info", "Set additional info.")
	printCmd.Arg(&value, true, "first_arg", "Value to be printed.")

	rootCmd := clic.NewFromFunc(func(ctx context.Context) error {
		fmt.Fprintln(out, "ouch, hit root")
		return nil
	}, "myapp", printCmd)

	// parse the cli command `myapp print --info=flagval arrrg`
	args := []string{"myapp", "print", "--info=flagval", "arrrg"}

	if err := rootCmd.Parse(args[1:]); err != nil {
		log.Fatalln(err)
	}

	if err := rootCmd.Handle(context.Background()); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// info flag = flagval
	// value arg = arrrg
}
