package clic

import (
	"fmt"
)

type Error struct {
	cause error
	c     *Clic
}

func NewError(err error, c *Clic) *Error {
	return &Error{err, c}
}

func (e *Error) Error() string {
	return fmt.Sprintf("cli command: %v", e.cause)
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Clic() *Clic {
	return e.c
}

type ParseError struct {
	cause error
	c     *Clic
}

func NewParseError(err error, c *Clic) *ParseError {
	return &ParseError{err, c}
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse: %v", e.cause)
}

func (e *ParseError) Unwrap() error {
	return e.cause
}

func (e *ParseError) Clic() *Clic {
	return e.c
}

type SubRequiredError struct {
	c *Clic
}

func NewSubRequiredError(c *Clic) *SubRequiredError {
	return &SubRequiredError{c}
}

func (e *SubRequiredError) Error() string {
	return fmt.Sprintf("a subcommand is required and not set")
}

func (e *SubRequiredError) Clic() *Clic {
	return e.c
}

type ArgSetError struct {
	cause error
}

func NewArgSetError(err error) *ArgSetError {
	return &ArgSetError{err}
}

func (e *ArgSetError) Error() string {
	return fmt.Sprintf("non-command args: %v", e.cause)
}

func (e *ArgSetError) Unwrap() error {
	return e.cause
}

type ArgMissingError struct {
	a *Arg
}

func NewArgMissingError(a *Arg) *ArgMissingError {
	return &ArgMissingError{a}
}

func (e *ArgMissingError) Error() string {
	return fmt.Sprintf("missing an expected arg: %s", e.a.Name)
}
