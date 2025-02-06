package clic_test

import (
	"context"

	"github.com/daved/clic"
)

func Example() {
	// error handling omitted to keep example focused

	var (
		info  = "default"
		value = "unset"
	)

	// Associate HandlerFunc with command name
	c := clic.NewFromFunc(printFunc(&info, &value), "myapp")

	// Associate flag and operand variables with relevant names
	c.Flag(&info, "i|info", "Set additional info.")
	c.Operand(&value, true, "first_operand", "Value to be printed.")

	// Parse the cli command as `myapp --info=flagval arrrg`
	_ = c.Parse([]string{"--info=flagval", "arrrg"})

	// Run the handler that Parse resolved to
	_ = c.HandleResolvedCmd(context.Background())

	// Output:
	// info flag = flagval
	// value operand = arrrg
}

func Example_aliases() {
	// error handling omitted to keep example focused

	// Associate HandlerFuncs with commands, setting "hello" as subcommand with an alias
	hc := clic.NewFromFunc(hello, "hello|aliased|h")
	c := clic.NewFromFunc(print, "myapp", hc)

	// Associate flag variables with relevant names; Technically unused here
	var debug bool
	c.Flag(&debug, "d|debug", "Set debug.")

	// Parse the cli command as `myapp -d aliased`
	_ = c.Parse([]string{"-d", "aliased"})

	// Run the handler that Parse resolved to
	_ = c.HandleResolvedCmd(context.Background())

	// Output:
	// Hello, World
}
