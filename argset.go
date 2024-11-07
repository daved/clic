package clic

import (
	"errors"
	"reflect"
)

type Arg struct {
	Val  any
	Opt  bool
	Name string
	Desc string
}

type ArgSet struct {
	as []*Arg
}

func NewArgSet() *ArgSet {
	return &ArgSet{}
}

func (as *ArgSet) Parse(args []string) error {
	for i, arg := range as.as {
		if len(args) <= i {
			return errors.New("not enough args") // TODO: check if args are opt, etc.
		}

		v := reflect.ValueOf(arg.Val)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if !v.CanSet() {
			return errors.New("cannot set argset val") // TODO: set this
		}

		v.SetString(args[i])
	}

	return nil
}

func (as *ArgSet) Arg(val any, opt bool, name, desc string) *ArgSet {
	as.as = append(as.as, &Arg{val, opt, name, desc})
	return as
}

// ArgSetProvider describes types that share an ArgSet.
type ArgSetProvider interface {
	ArgSet() *ArgSet
}
