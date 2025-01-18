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

	// Associate "print" flag and operand variables with relevant names
	c.Flag(&info, "i|info", "Set additional info.")
	c.Operand(&value, true, "first_operand", "Value to be printed.")

	// Parse the cli command as `myapp --info=flagval arrrg`
	_ = c.Parse([]string{"--info=flagval", "arrrg"})

	// Run the handler that Parse resolved to
	_ = c.HandleResolvedCmd(context.Background())

	// Output:
	// info flag = flagval
	// value arg = arrrg
}
