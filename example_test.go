package clic_test

import (
	"context"
	"fmt"
	"os"

	"github.com/daved/clic"
)

func Example() {
	cmd := NewRootClic("myapp",
		NewSubClic("subcmd"),
	)

	// parse the cli command `myapp subcmd --info=flagval`
	args := []string{"myapp", "subcmd", "--info=flagval"}

	if err := cmd.Parse(args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := cmd.Handle(context.Background()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Output:
	// info = flagval
}

type Root struct{}

func NewRootClic(name string, subs ...*clic.Clic) *clic.Clic {
	return clic.New(&Root{}, name, subs...)
}

func (cmd *Root) HandleCommand(ctx context.Context) error {
	fmt.Println("hit root")
	return nil
}

type Sub struct {
	info string
}

func NewSubClic(name string) *clic.Clic {
	cmd := &Sub{info: "default"}

	c := clic.New(cmd, name)
	c.Flag(&cmd.info, "i|info", "set info value")

	return c
}

func (cmd *Sub) HandleCommand(ctx context.Context) error {
	fmt.Printf("info = %s\n", cmd.info)
	return nil
}
