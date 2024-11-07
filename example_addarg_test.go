package clic_test

import (
	"context"
	"fmt"
	"os"

	"github.com/daved/clic"
	"github.com/daved/flagset"
)

func ExampleClic_AddArg() {
	subCmd := clic.New(NewPrintFirstArgCmd("printarg"))
	rootCmd := clic.New(NewRootCmd(), subCmd)

	// parse the cli command `myapp printarg --info=flagval arrrg`
	if err := rootCmd.Parse([]string{"printarg", "--info=flagval", "arrrg"}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := rootCmd.Handle(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Output:
	// flagval
	// First Arg = arrrg
}

type PrintFirstArgCmd struct {
	fs   *flagset.FlagSet
	as   *clic.ArgSet
	Info string
	Arg  string
}

func NewPrintFirstArgCmd(name string) *PrintFirstArgCmd {
	fs := flagset.New(name)
	as := clic.NewArgSet()

	cmd := &PrintFirstArgCmd{
		fs:   fs,
		as:   as,
		Info: "default",
		Arg:  "unset",
	}

	fs.Opt(&cmd.Info, "i|info", "set info value")
	as.Arg(&cmd.Arg, true, "first_arg", "First arg, will be printed.")

	return cmd
}

func (cmd *PrintFirstArgCmd) FlagSet() *flagset.FlagSet {
	return cmd.fs
}

func (cmd *PrintFirstArgCmd) ArgSet() *clic.ArgSet {
	return cmd.as
}

func (cmd *PrintFirstArgCmd) HandleCommand(ctx context.Context) error {
	fmt.Printf("%s\nFirst Arg = %v\n", cmd.Info, cmd.Arg)
	return nil
}
