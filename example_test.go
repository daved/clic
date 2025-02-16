package clic_test

import (
	"context"
	"fmt"

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
	resolved, _ := c.Parse([]string{"--info=flagval", "arrrg"})

	// Run the handler that Parse resolved to
	_ = resolved.Handle(context.Background())

	// Output:
	// info flag = flagval
	// value operand = arrrg
}

func Example_aliases() {
	// error handling omitted to keep example focused

	// Associate HandlerFuncs with commands, setting "hello" as subcommand with an alias
	hc := clic.NewFromFunc(hello, "hello|aliased")
	c := clic.NewFromFunc(print, "myapp", hc)

	// Associate flag variables with relevant names; Technically unused here
	var debug bool
	c.Flag(&debug, "d|debug", "Set debug.")

	// Parse the cli command as `myapp -d aliased`
	resolved, _ := c.Parse([]string{"-d", "aliased"})

	// Run the handler that Parse resolved to
	_ = resolved.Handle(context.Background())

	// Output:
	// Hello, World
}

func Example_categories() {
	// error handling omitted to keep example focused

	// Associate HandlerFuncs with commands, setting cat and desc fields
	hc1 := clic.NewFromFunc(hello, "hello1")
	hc1.Category = "Foo"
	hc1.Description = "Say hello Uno"

	hc2 := clic.NewFromFunc(hello, "hello2")
	hc2.Category = "Foo"
	hc2.Description = "Hello hello xoxo"

	hc := clic.NewFromFunc(hello, "hello")
	hc.Category = "Bar"
	hc.Description = "Helo 0000 00 000"

	// Control subcommand category order in the parent
	c := clic.NewFromFunc(print, "myapp", hc1, hc2, hc)
	c.SubCmdCatsSort = []string{"Foo|Foo-related", "Bar|All things Bar"}

	// Parse the cli command as `myapp`
	resolved, _ := c.Parse([]string{})

	fmt.Println(resolved.Usage())

	// Output:
	// Usage:
	//
	//   myapp [hello1|hello2|hello]
	//
	// Subcommands for myapp:
	//
	//   Foo          Foo-related
	//     hello1         Say hello Uno
	//     hello2         Hello hello xoxo
	//
	//   Bar          All things Bar
	//     hello          Helo 0000 00 000
}
