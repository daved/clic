package clic_test

import (
	"context"
	"errors"
	"fmt"
	"os"
)

func printFunc(info, value *string) func(context.Context) error {
	return func(ctx context.Context) error {
		fmt.Printf("info flag = %s\nvalue operand = %v\n", *info, *value)
		return nil
	}
}

func print(ctx context.Context) error {
	var (
		info  = "default"
		value = "unset"
	)

	return printFunc(&info, &value)(ctx)
}

func hello(context.Context) error {
	fmt.Println("Hello, World")
	return nil
}

func details(context.Context) error {
	fmt.Printf("Args: %v\n", os.Args)
	return nil
}

func unused(context.Context) error {
	return errors.New("unused")
}
