package clic_test

import (
	"context"
	"fmt"
	"os"

	"github.com/daved/clic"
)

func ExampleClic_ArgSet() {
	rootCmd := NewRootClic("myapp",
		NewPrintFirstArgClic("printarg"),
	)

	// parse the cli command `myapp printarg --info=flagval arrrg`
	args := []string{"myapp", "printarg", "--info=flagval", "arrrg"}

	if err := rootCmd.Parse(args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := rootCmd.Handle(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Output:
	// info = flagval
	// first arg = arrrg
}

type PrintFirstArg struct {
	info string
	arg  string
}

func NewPrintFirstArgClic(name string) *clic.Clic {
	cmd := &PrintFirstArg{
		info: "default",
		arg:  "unset",
	}

	c := clic.New(cmd, name)

	c.FlagSet.Opt(&cmd.info, "i|info", "set info value")
	c.ArgSet.Arg(&cmd.arg, true, "first_arg", "First arg, will be printed.")

	return c
}

func (cmd *PrintFirstArg) HandleCommand(ctx context.Context) error {
	fmt.Printf("info = %s\n", cmd.info)
	fmt.Printf("first arg = %v\n", cmd.arg)
	return nil
}
