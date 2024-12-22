package argset

import (
	"errors"
	"reflect"

	errs "github.com/daved/clic/clicerrs"
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

func New() *ArgSet {
	return &ArgSet{}
}

func (as *ArgSet) Parse(args []string) error {
	for i, arg := range as.Args {
		if len(args) <= i {
			if !arg.Req {
				continue
			}

			return errs.NewArgSetError(errs.NewArgMissingError(arg.Name))
		}

		v := reflect.ValueOf(arg.Val)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if !v.CanSet() {
			return errs.NewArgSetError(errors.New("unsettable value used for arg"))
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
