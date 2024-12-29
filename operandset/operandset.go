package operandset

import (
	"flag"
	"strconv"
	"time"

	errs "github.com/daved/clic/operandset/oserrs"
	"github.com/daved/clic/operandset/vtypes"
)

type Operand struct {
	val  any
	req  bool
	name string
	desc string
	Tag  string
	Meta map[string]any
}

func (o *Operand) Name() string {
	return o.name
}

func (o *Operand) Required() bool {
	return o.req
}

func (o *Operand) Description() string {
	return o.desc
}

type OperandSet struct {
	name    string
	ops     []*Operand
	raws    []string
	tmplCfg *TmplConfig
	Meta    map[string]any
}

func New(name string) *OperandSet {
	return &OperandSet{
		name:    name,
		tmplCfg: NewDefaultTmplConfig(),
		Meta:    map[string]any{},
	}
}

func (os *OperandSet) Name() string {
	return os.name
}

func (os *OperandSet) Operands() []*Operand {
	return os.ops
}

func (os *OperandSet) Operand(val any, req bool, name, desc string) *Operand {
	lEnc, rEnc := "[", "]" // enclosures
	if req {
		lEnc, rEnc = "<", ">"
	}

	o := &Operand{
		val:  val,
		req:  req,
		name: name,
		desc: desc,
		Tag:  lEnc + name + rEnc,
		Meta: map[string]any{},
	}

	os.ops = append(os.ops, o)

	return o
}

func (os *OperandSet) Parse(args []string) error {
	os.raws = args

	if err := parse(os.ops, args); err != nil {
		return NewError(err)
	}

	return nil
}

func parse(ops []*Operand, args []string) error {
	newError := errs.NewParseError

	for i, op := range ops {
		if len(args) <= i {
			if !op.req {
				continue
			}

			return newError(errs.NewOperandMissingError(op.name))
		}

		raw := args[i]

		switch v := op.val.(type) {
		case *string:
			*v = raw

		case *bool:
			b, err := strconv.ParseBool(raw)
			if err != nil {
				return newError(errs.NewConvertRawError(err))
			}
			*v = b

		case *int:
			n, err := strconv.Atoi(raw)
			if err != nil {
				return newError(errs.NewConvertRawError(err))
			}
			*v = n

		case *int64:
			n, err := strconv.ParseInt(raw, 10, 0)
			if err != nil {
				return newError(errs.NewConvertRawError(err))
			}
			*v = n

		case *uint:
			n, err := strconv.ParseUint(raw, 10, 0)
			if err != nil {
				return newError(errs.NewConvertRawError(err))
			}
			*v = uint(n)

		case *uint64:
			n, err := strconv.ParseUint(raw, 10, 0)
			if err != nil {
				return newError(errs.NewConvertRawError(err))
			}
			*v = n

		case *float64:
			f, err := strconv.ParseFloat(raw, 64)
			if err != nil {
				return newError(errs.NewConvertRawError(err))
			}
			*v = f

		case *time.Duration:
			d, err := time.ParseDuration(raw)
			if err != nil {
				return newError(errs.NewConvertRawError(err))
			}
			*v = d

		case vtypes.TextMarshalUnmarshaler:
			if err := v.UnmarshalText([]byte(raw)); err != nil {
				return newError(errs.NewConvertRawError(err))
			}

		case flag.Value:
			if err := v.Set(raw); err != nil {
				return newError(errs.NewConvertRawError(err))
			}

		case vtypes.OperandFunc:
			if err := v(raw); err != nil {
				return newError(errs.NewConvertRawError(err))
			}
		}
	}

	return nil
}

func (os *OperandSet) Parsed() []string {
	return os.raws
}

func (os *OperandSet) SetUsageTemplating(tmplCfg *TmplConfig) {
	os.tmplCfg = tmplCfg
}

func (os *OperandSet) Usage() string {
	return executeTmpl(os.tmplCfg, &TmplData{OperandSet: os})
}
