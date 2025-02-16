package clic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/daved/clic/operandset"
)

func defaultPtrs[T any](args ...T) []any {
	var ptrs []any
	for _, arg := range args {
		a := arg // must copy
		ptrs = append(ptrs, &a)
	}
	return ptrs
}

func TestClicParse(t *testing.T) {
	type parseScope struct {
		name   string
		clicFn func(*bytes.Buffer, *[]any) *Clic
	}

	scopeA := parseScope{
		name: "clic subcmd-opt opnd0-req opnd1-opt",
		clicFn: func(buf *bytes.Buffer, ptrs *[]any) *Clic {
			return NewCmdClic(buf, "myapp", nil,
				NewCmdClic(buf, "subcmd",
					func(os *operandset.OperandSet) {
						*ptrs = defaultPtrs("default0", "default1")
						os.Operand((*ptrs)[0], true, "first_opnd", "")
						os.Operand((*ptrs)[1], false, "second_opnd", "")
					},
				),
			)
		},
	}

	scopeB := parseScope{
		name: "clic subcmd-req opnd0-req opnd1-opt",
		clicFn: func(buf *bytes.Buffer, ptrs *[]any) *Clic {
			c := NewCmdClic(buf, "myapp", nil,
				NewCmdClic(buf, "subcmd",
					func(os *operandset.OperandSet) {
						*ptrs = defaultPtrs("default0", "default1")
						os.Operand((*ptrs)[0], true, "first_opnd", "")
						os.Operand((*ptrs)[1], false, "second_opnd", "")
					},
				),
			)
			c.SubRequired = true
			return c
		},
	}

	tt := []struct {
		scope parseScope
		name  string
		args  []string
		out   string
		vals  []any
		cause error
	}{
		{
			scope: scopeA,
			name:  "subcmd-one opnds-both",
			args: []string{
				"myapp", "subcmd", "--info=flagval",
				"first", "second",
			},
			out:  "subcmd",
			vals: []any{"first", "second"},
		},
		{
			scope: scopeA,
			name:  "subcmd-one opnds-none",
			args: []string{
				"myapp", "subcmd", "--info=flagval",
			},
			cause: CauseOperandRequired,
		},
		{
			scope: scopeA,
			name:  "subcmd-one opnds-first",
			args: []string{
				"myapp", "subcmd", "--info=flagval",
				"first",
			},
			out:  "subcmd",
			vals: []any{"first", "default1"},
		},
		{
			scope: scopeB,
			name:  "subcmd-one opnds-both",
			args: []string{
				"myapp", "subcmd", "--info=flagval",
				"first", "second",
			},
			out:  "subcmd",
			vals: []any{"first", "second"},
		},
		{
			scope: scopeB,
			name:  "subcmd-none opnds-both",
			args: []string{
				"myapp", "--info=flagval",
				"first", "second",
			},
			cause: CauseSubCmdRequired,
		},
		{
			scope: scopeB,
			name:  "subcmd-none opnds-none",
			args: []string{
				"myapp",
			},
			cause: CauseSubCmdRequired,
		},
	}

	for _, tc := range tt {
		t.Run(tc.scope.name+"/"+tc.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			var ptrs []any

			c := tc.scope.clicFn(buf, &ptrs)

			err := c.Parse(tc.args[1:])
			if !errors.Is(err, tc.cause) {
				t.Fatalf("parse error: got %v, want %v", err, tc.cause)
			}
			if err != nil {
				return
			}

			_ = c.HandleResolvedCmd(context.Background())

			out := buf.String()
			if out != tc.out {
				t.Fatalf("output: got: %v, want: %v", out, tc.out)
			}

			for i, ptr := range ptrs {
				got := reflect.ValueOf(ptr).Elem()
				want := reflect.ValueOf(tc.vals[i])
				if !got.Equal(want) {
					t.Fatalf("vals: got: %v, want: %v", got, want)
				}
			}
		})
	}
}

type Cmd struct {
	buf  *bytes.Buffer
	name string
}

type setupFunc func(*operandset.OperandSet)

func NewCmdClic(buf *bytes.Buffer, name string, fn setupFunc, subs ...*Clic) *Clic {
	cmd := &Cmd{
		buf:  buf,
		name: name,
	}

	var info string
	var num int

	c := New(cmd, name, subs...)
	c.Flag(&info, "info|i", "")
	c.Flag(&num, "num|n", "")

	if fn != nil {
		fn(c.OperandSet)
	}

	return c
}

func (cmd *Cmd) HandleCommand(ctx context.Context) error {
	cmd.buf.Reset()
	fmt.Fprintf(cmd.buf, "%s", cmd.name)
	return nil
}
