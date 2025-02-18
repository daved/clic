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
	root := clic.NewFromFunc(printFunc(&info, &value), "myapp")

	// Associate flag and operand variables with relevant names
	root.Flag(&info, "i|info", "Set additional info.")
	root.Operand(&value, true, "first_operand", "Value to be printed.")

	// Parse the cli command as `myapp --info=flagval arrrg`
	cmd, _ := root.Parse([]string{"--info=flagval", "arrrg"})

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
	root := clic.NewFromFunc(hello, "hello|aliased")

	// Parse the cli command as `myapp aliased`
	cmd, _ := root.Parse([]string{"aliased"})

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

	details := clic.NewFromFunc(details, "details")
	details.Category = "Informational"
	details.Description = "List details (os.Args)"

	// Associate HandlerFunc with command name
	root := clic.NewFromFunc(unused, "myapp", hello, goodbye, details)
	root.SubRequired = true
	// Set up subcommand category order
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
	root := clic.NewFromFunc(hello, "myapp")

	// Associate flag variable with relevant name
	root.Flag(&verbosity, "v", "Set verbosity. Can be set multiple times.")

	// Parse the cli command as `myapp -vvv`
	cmd, _ := root.Parse([]string{"-vvv"})

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

func Example_handlerWrapping() {
	// error handling omitted to keep example focused

	// Associate HandlerFuncs with command names
	hello := clic.NewFromFunc(hello, "hello")
	goodbye := clic.NewFromFunc(goodbye, "goodbye")
	details := clic.NewFromFunc(details, "details")

	root := clic.NewFromFunc(unused, "myapp", hello, goodbye, details)
	root.SubRequired = true

	root.Recursively(func(c *clic.Clic) {
		next := c.Handler.HandleCommand
		c.Handler = clic.HandlerFunc(func(ctx context.Context) error {
			fmt.Println("before")

			if err := next(ctx); err != nil {
				return err
			}

			fmt.Println("after")
			return nil
		})
	})

	// Parse the cli command as `myapp hello`
	cmd, _ := root.Parse([]string{"hello"})

	// Run the handler that Parse resolved to
	_ = cmd.Handle(context.Background())
	// Output:
	// before
	// Hello, World
	// after
}
