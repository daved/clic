package clic_test

import (
	"context"
	"fmt"
	"os"

	"github.com/daved/clic"
	"github.com/daved/flagset"
)

func Example() {
	subCmd := clic.New(NewSubCmd("subcmd"))
	rootCmd := clic.New(NewRootCmd(), subCmd)

	// parse the cli command `myapp subcmd --info=flagval arrrg`
	if err := rootCmd.Parse([]string{"subcmd", "--info=flagval", "arrrg"}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := rootCmd.Handle(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Output:
	// flagval
	// [arrrg]
}

type RootCmd struct {
	fs *flagset.FlagSet
}

func NewRootCmd() *RootCmd {
	return &RootCmd{fs: flagset.New("root")}
}

func (cmd *RootCmd) FlagSet() *flagset.FlagSet {
	return cmd.fs
}

func (cmd *RootCmd) HandleCommand(ctx context.Context) error {
	fmt.Println("hit root")
	return nil
}

type SubCmd struct {
	fs   *flagset.FlagSet
	info string
}

func NewSubCmd(name string) *SubCmd {
	fs := flagset.New(name)

	cmd := &SubCmd{
		fs:   fs,
		info: "default",
	}

	fs.Opt(&cmd.info, "i|info", "set info value")

	return cmd
}

func (cmd *SubCmd) FlagSet() *flagset.FlagSet {
	return cmd.fs
}

func (cmd *SubCmd) HandleCommand(ctx context.Context) error {
	fmt.Printf("%s\n%v\n", cmd.info, cmd.fs.Args())
	return nil
}
