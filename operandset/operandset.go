// Package operandset wraps the [operandset] package.
package operandset

import "github.com/daved/operandset"

// Operand is an alias of [operandset.Operand].
type Operand = operandset.Operand

// OperandSet wraps [operandset.OperandSet].
type OperandSet struct {
	*operandset.OperandSet
}

// New returns an instance of OperandSet.
func New(name string) *OperandSet {
	return &OperandSet{
		OperandSet: operandset.New(name),
	}
}
