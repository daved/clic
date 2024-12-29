package oserrs

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
	return fmt.Sprintf("convert raw string: %v", e.child)
}

func (e *ParseError) Unwrap() error {
	return e.child
}

func (e *ParseError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type OperandMissingError struct {
	name string
}

func NewOperandMissingError(name string) *OperandMissingError {
	return &OperandMissingError{name}
}

func (e *OperandMissingError) Error() string {
	return fmt.Sprintf("missing an expected operand: %s", e.name)
}

func (e *OperandMissingError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}

type ConvertRawError struct {
	child error
}

func NewConvertRawError(child error) *ConvertRawError {
	return &ConvertRawError{child}
}

func (e *ConvertRawError) Error() string {
	return fmt.Sprintf("convert raw string: %v", e.child)
}

func (e *ConvertRawError) Unwrap() error {
	return e.child
}

func (e *ConvertRawError) Is(err error) bool {
	return reflect.TypeOf(e) == reflect.TypeOf(err)
}
