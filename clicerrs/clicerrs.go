package clicerrs

import (
	"fmt"
	"reflect"
)

type ParseError struct {
	child error
}

func NewParseError(child error) *ParseError {
	return &ParseError{child}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse: %v", e.child)
}

func (e *ParseError) Unwrap() error {
	return e.child
}

func (e *ParseError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type SubRequiredError struct{}

func NewSubRequiredError() *SubRequiredError {
	return &SubRequiredError{}
}

func (e *SubRequiredError) Error() string {
	return "subcommand is required and not set"
}

func (e *SubRequiredError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type FlagSetError struct {
	child error
}

func NewFlagSetError(child error) *FlagSetError {
	return &FlagSetError{child}
}

func (e *FlagSetError) Error() string {
	return e.child.Error()
}

func (e *FlagSetError) Unwrap() error {
	return e.child
}

func (e *FlagSetError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}
