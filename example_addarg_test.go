package clic_test

import (
	"context"
	"fmt"
	"os"

	"github.com/daved/clic"
)

func ExampleClic_Arg() {
	subCmdHandler := NewPrintFirstArg()
	subCmd := clic.New(subCmdHandler, "printarg")

	subCmd.Flag(&subCmdHandler.info, "i|info", "set info value")
	subCmd.Arg(&subCmdHandler.arg, true, "first_arg", "First arg, will be printed.")

	rootCmdHandler := &Root{}
	rootCmd := clic.New(rootCmdHandler, "myapp", subCmd)

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

type Root struct{}

func (cmd *Root) HandleCommand(ctx context.Context) error {
	fmt.Println("hit root")
	return nil
}

type PrintFirstArg struct {
	info string
	arg  string
}

func NewPrintFirstArg() *PrintFirstArg {
	return &PrintFirstArg{
		info: "default",
		arg:  "unset",
	}
}

func (cmd *PrintFirstArg) HandleCommand(ctx context.Context) error {
	fmt.Printf("info = %s\n", cmd.info)
	fmt.Printf("first arg = %v\n", cmd.arg)
	return nil
}
