package clic_test

import (
	"context"
	"fmt"
	"os"
)

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
