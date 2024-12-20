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

type ArgSetError struct {
	child error
}

func NewArgSetError(child error) *ArgSetError {
	return &ArgSetError{child}
}

func (e *ArgSetError) Error() string {
	return fmt.Sprintf("non-command args: %v", e.child)
}

func (e *ArgSetError) Unwrap() error {
	return e.child
}

func (e *ArgSetError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type ArgMissingError struct {
	name string
}

func NewArgMissingError(name string) *ArgMissingError {
	return &ArgMissingError{name}
}

func (e *ArgMissingError) Error() string {
	return fmt.Sprintf("missing an expected arg: %s", e.name)
}

func (e *ArgMissingError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}
