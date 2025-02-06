package clic_test

import (
	"context"
	"fmt"
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
