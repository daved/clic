package operandset

import "github.com/daved/operandset"

type Operand = operandset.Operand

type OperandSet struct {
	*operandset.OperandSet
}

func New(name string) *OperandSet {
	return &OperandSet{
		OperandSet: operandset.New(name),
	}
}
