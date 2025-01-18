package clic_test

import (
	"context"
	"fmt"
)

func printFunc(info, value *string) func(context.Context) error {
	return func(ctx context.Context) error {
		fmt.Printf("info flag = %s\nvalue arg = %v\n", *info, *value)
		return nil
	}
}
