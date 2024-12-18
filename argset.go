package clic

import (
	"errors"
	"reflect"
)

type Arg struct {
	Val  any
	Req  bool
	Name string
	Desc string
	Meta map[string]any
	Hint string
}

type ArgSet struct {
	Args []*Arg
}

func newArgSet() *ArgSet {
	return &ArgSet{}
}

func (as *ArgSet) parse(args []string) error {
	for i, arg := range as.Args {
		if len(args) <= i {
			if !arg.Req {
				continue
			}

			return NewArgSetError(NewArgMissingError(arg))
		}

		v := reflect.ValueOf(arg.Val)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if !v.CanSet() {
			return NewArgSetError(errors.New("unsettable value used for arg"))
		}

		v.SetString(args[i])
		// TODO: handle other types
	}

	return nil
}

func (as *ArgSet) Arg(val any, req bool, name, desc string) *Arg {
	a := &Arg{
		Val:  val,
		Req:  req,
		Name: name,
		Desc: desc,
		Meta: make(map[string]any),
	}

	as.Args = append(as.Args, a)

	lEnc, rEnc := "[", "]" // enclosures
	if a.Req {
		lEnc, rEnc = "<", ">"
	}

	a.Hint = lEnc + a.Name + rEnc

	return a
}
