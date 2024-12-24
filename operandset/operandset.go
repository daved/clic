package operandset

import (
	"encoding"
	"flag"
	"strconv"
	"time"

	errs "github.com/daved/clic/clicerrs"
)

type TextMarshalUnmarshaler interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
}

type OperandFunc func(string) error

type Operand struct {
	Val  any
	Req  bool
	Name string
	Desc string
	Tag  string
	Meta map[string]any
}

type OperandSet struct {
	Operands []*Operand
}

func New() *OperandSet {
	return &OperandSet{}
}

func (os *OperandSet) Parse(args []string) error {
	for i, op := range os.Operands {
		if len(args) <= i {
			if !op.Req {
				continue
			}

			return errs.NewOperandSetError(errs.NewOperandMissingError(op.Name))
		}

		raw := args[i]

		switch v := op.Val.(type) {
		case *string:
			*v = raw

		case *bool:
			b, err := strconv.ParseBool(raw)
			if err != nil {
				return errs.NewOperandSetError(err)
			}
			*v = b

		case *int:
			n, err := strconv.Atoi(raw)
			if err != nil {
				return errs.NewOperandSetError(err)
			}
			*v = n

		case *int64:
			n, err := strconv.ParseInt(raw, 10, 0)
			if err != nil {
				return errs.NewOperandSetError(err)
			}
			*v = n

		case *uint:
			n, err := strconv.ParseUint(raw, 10, 0)
			if err != nil {
				return errs.NewOperandSetError(err)
			}
			*v = uint(n)

		case *uint64:
			n, err := strconv.ParseUint(raw, 10, 0)
			if err != nil {
				return errs.NewOperandSetError(err)
			}
			*v = n

		case *float64:
			f, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return errs.NewOperandSetError(err)
			}
			*v = f

		case *time.Duration:
			d, err := time.ParseDuration(raw)
			if err != nil {
				return errs.NewOperandSetError(err)
			}
			*v = d

		case TextMarshalUnmarshaler:
			if err := v.UnmarshalText([]byte(raw)); err != nil {
				return err
			}

		case flag.Value:
			if err := v.Set(raw); err != nil {
				return err
			}

		case OperandFunc:
			if err := v(raw); err != nil {
				return err
			}
		}
	}

	return nil
}

func (os *OperandSet) Operand(val any, req bool, name, desc string) *Operand {
	o := &Operand{
		Val:  val,
		Req:  req,
		Name: name,
		Desc: desc,
		Meta: map[string]any{},
	}

	os.Operands = append(os.Operands, o)

	lEnc, rEnc := "[", "]" // enclosures
	if o.Req {
		lEnc, rEnc = "<", ">"
	}

	o.Tag = lEnc + o.Name + rEnc

	return o
}
