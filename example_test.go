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
	cmd, _ := c.Parse([]string{"--info=flagval", "arrrg"})

	// Run the handler that Parse resolved to
	_ = cmd.Handle(context.Background())

	fmt.Println()
	fmt.Println(cmd.Usage())
	// Output:
	// info flag = flagval
	// value operand = arrrg
	//
	// Usage:
	//
	//   myapp [FLAGS] <first_operand>
	//
	// Flags for myapp:
	//
	//     -i, --info  =STRING    default: default
	//         Set additional info.
}

func Example_aliases() {
	// error handling omitted to keep example focused

	// Associate HandlerFunc with command name and alias
	c := clic.NewFromFunc(hello, "hello|aliased")

	// Parse the cli command as `myapp aliased`
	cmd, _ := c.Parse([]string{"aliased"})

	// Run the handler that Parse resolved to
	_ = cmd.Handle(context.Background())

	// Output:
	// Hello, World
}

func Example_categories() {
	// Associate HandlerFuncs with command names, setting cat and desc fields
	hello := clic.NewFromFunc(hello, "hello")
	hello.Category = "Salutations"
	hello.Description = "Show hello world message"

	goodbye := clic.NewFromFunc(goodbye, "goodbye")
	goodbye.Category = "Salutations"
	goodbye.Description = "Show goodbye message"

	print := clic.NewFromFunc(details, "details")
	print.Category = "Informational"
	print.Description = "List details (os.Args)"

	// Set up subcommand category order in the parent
	root := clic.NewFromFunc(unused, "myapp", hello, goodbye, print)
	root.SubRequired = true
	// Category names seperated from optional descriptions by "|"
	root.SubCmdCatsSort = []string{"Salutations|Salutations-related", "Informational|All things info"}

	// Parse the cli command as `myapp`; will return error from lack of subcommand
	cmd, err := root.Parse([]string{})
	if err != nil {
		fmt.Println(cmd.Usage())
		fmt.Println(err)
	}
	// Output:
	// Usage:
	//
	//   myapp {hello|goodbye|details}
	//
	// Subcommands for myapp:
	//
	//   Salutations        Salutations-related
	//     hello              Show hello world message
	//     goodbye            Show goodbye message
	//
	//   Informational      All things info
	//     details            List details (os.Args)
	//
	// cli command: parse: subcommand required
}

func Example_verbosity() {
	// error handling omitted to keep example focused

	var verbosity []bool

	// Associate HandlerFunc with command name
	c := clic.NewFromFunc(hello, "myapp")

	// Associate flag variable with relevant name
	c.Flag(&verbosity, "v", "Set verbosity. Can be set multiple times.")

	// Parse the cli command as `myapp -vvv`
	cmd, _ := c.Parse([]string{"-vvv"})

	// Run the handler that Parse resolved to
	_ = cmd.Handle(context.Background())

	fmt.Printf("verbosity: length=%d value=%v\n", len(verbosity), verbosity)
	fmt.Println()
	fmt.Println(cmd.Usage())
	// Output:
	// Hello, World
	// verbosity: length=3 value=[true true true]
	//
	// Usage:
	//
	//   myapp [FLAGS]
	//
	// Flags for myapp:
	//
	//     -v  =BOOL
	//         Set verbosity. Can be set multiple times.
}
