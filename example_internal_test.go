package clic_test

import (
	"context"
	"fmt"
	"os"
)

func printFunc(info, value *string) func(context.Context) error {
	return func(ctx context.Context) error {
		fmt.Printf("info flag = %s\nvalue operand = %v\n", *info, *value)
		return nil
	}
}

func hello(context.Context) error {
	fmt.Println("Hello, World")
	return nil
}

func goodbye(context.Context) error {
	fmt.Println("Goodbye")
	return nil
}

func details(context.Context) error {
	fmt.Printf("Args: %v\n", os.Args)
	return nil
}

func printRoot(context.Context) error {
	fmt.Println("Root")
	return nil
}
