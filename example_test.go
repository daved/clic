package clic_test

import (
	"context"
	"fmt"

	"github.com/daved/clic"
	"github.com/daved/flagset"
)

func Example() {
	subCmd := clic.New(NewSubCmd("sub"))
	rootCmd := clic.New(NewRootCmd("root"), subCmd)

	if err := rootCmd.Parse([]string{"sub", "--info=flag", "arrrg"}); err != nil {
		fmt.Println(err)
		return
	}

	if err := rootCmd.HandleCalled(context.Background()); err != nil {
		fmt.Println(err)
		return
	}

	// Output:
	// flag
	// [arrrg]
}

type RootCmd struct {
	fs *flagset.FlagSet
}

func NewRootCmd(name string) *RootCmd {
	return &RootCmd{
		fs: flagset.New(name),
	}
}

func (cmd *RootCmd) FlagSet() *flagset.FlagSet {
	return cmd.fs
}

func (cmd *RootCmd) HandleCommand(ctx context.Context, c *clic.Clic) error {
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

func (cmd *SubCmd) HandleCommand(ctx context.Context, c *clic.Clic) error {
	fmt.Println(cmd.info)
	fmt.Println(c.FlagSet().Args())
	return nil
}
