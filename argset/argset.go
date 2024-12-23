package argset

import (
	"encoding"
	"strconv"
	"time"

	errs "github.com/daved/clic/clicerrs"
)

type TextMarshalUnmarshaler interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

type ArgValue interface {
	String()
	Set(string) error
}

type ArgFunc func(string) error

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

		raw := args[i]

		switch v := arg.Val.(type) {
		case *string:
			*v = raw

		case *bool:
			b, err := strconv.ParseBool(raw)
			if err != nil {
				return errs.NewArgSetError(err)
			}
			*v = b

		case *int:
			n, err := strconv.Atoi(raw)
			if err != nil {
				return errs.NewArgSetError(err)
			}
			*v = n

		case *int64:
			n, err := strconv.ParseInt(raw, 10, 0)
			if err != nil {
				return errs.NewArgSetError(err)
			}
			*v = n

		case *uint:
			n, err := strconv.ParseUint(raw, 10, 0)
			if err != nil {
				return errs.NewArgSetError(err)
			}
			*v = uint(n)

		case *uint64:
			n, err := strconv.ParseUint(raw, 10, 0)
			if err != nil {
				return errs.NewArgSetError(err)
			}
			*v = n

		case *float64:
			f, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return errs.NewArgSetError(err)
			}
			*v = f

		case *time.Duration:
			d, err := time.ParseDuration(raw)
			if err != nil {
				return errs.NewArgSetError(err)
			}
			*v = d

		case TextMarshalUnmarshaler:
			if err := v.UnmarshalText([]byte(raw)); err != nil {
				return err
			}

		case ArgValue:
			if err := v.Set(raw); err != nil {
				return err
			}

		case ArgFunc:
			if err := v(raw); err != nil {
				return err
			}
		}
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
