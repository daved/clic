package clic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
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
		name: "clic subcmd-opt arg0-req arg1-opt",
		clicFn: func(buf *bytes.Buffer, ptrs *[]any) *Clic {
			return NewCmdClic(buf, "myapp", nil,
				NewCmdClic(buf, "subcmd",
					func(as *ArgSet) {
						*ptrs = defaultPtrs("default0", "default1")
						as.Arg((*ptrs)[0], true, "first_arg", "")
						as.Arg((*ptrs)[1], false, "second_arg", "")
					},
				),
			)
		},
	}

	scopeB := parseScope{
		name: "clic subcmd-req arg0-req arg1-opt",
		clicFn: func(buf *bytes.Buffer, ptrs *[]any) *Clic {
			c := NewCmdClic(buf, "myapp", nil,
				NewCmdClic(buf, "subcmd",
					func(as *ArgSet) {
						*ptrs = defaultPtrs("default0", "default1")
						as.Arg((*ptrs)[0], true, "first_arg", "")
						as.Arg((*ptrs)[1], false, "second_arg", "")
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
			name:  "subcmd one args both",
			args: []string{
				"myapp", "subcmd", "--info=flagval",
				"first", "second",
			},
			out:  "subcmd",
			vals: []any{"first", "second"},
		},
		{
			scope: scopeA,
			name:  "subcmd one args none",
			args: []string{
				"myapp", "subcmd", "--info=flagval",
			},
			cause: CauseParseArgMissing,
		},
		{
			scope: scopeA,
			name:  "subcmd one args first",
			args: []string{
				"myapp", "subcmd", "--info=flagval",
				"first",
			},
			out:  "subcmd",
			vals: []any{"first", "default1"},
		},
		{
			scope: scopeB,
			name:  "subcmd one args both",
			args: []string{
				"myapp", "subcmd", "--info=flagval",
				"first", "second",
			},
			out:  "subcmd",
			vals: []any{"first", "second"},
		},
		{
			scope: scopeB,
			name:  "subcmd none args both",
			args: []string{
				"myapp", "--info=flagval",
				"first", "second",
			},
			cause: CauseParseSubRequired,
		},
		{
			scope: scopeB,
			name:  "subcmd none args none",
			args: []string{
				"myapp",
			},
			cause: CauseParseSubRequired,
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

			_ = c.Handle(context.Background())

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

type setupFunc func(*ArgSet)

func NewCmdClic(buf *bytes.Buffer, name string, fn setupFunc, subs ...*Clic) *Clic {
	cmd := &Cmd{
		buf:  buf,
		name: name,
	}

	var info string
	var num int

	c := New(cmd, name, subs...)
	c.FlagSet.FlagSet.Opt(&info, "info|i", "")
	c.FlagSet.FlagSet.Opt(&num, "num|n", "")

	if fn != nil {
		fn(c.ArgSet)
	}

	return c
}

func (cmd *Cmd) HandleCommand(ctx context.Context) error {
	cmd.buf.Reset()
	fmt.Fprintf(cmd.buf, "%s", cmd.name)
	return nil
}
