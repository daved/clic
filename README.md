# clic [![GoDoc](https://pkg.go.dev/badge/github.com/daved/clic.svg)](https://pkg.go.dev/github.com/daved/clic)

```go
go get github.com/daved/clic
```

Package clic provides a structured multiplexer for CLI commands. In other words, clic will parse CLI
command arguments and route callers to the appropriate handler.

## Usage

```go
type Clic
    func New(h Handler, name string, subs ...*Clic) *Clic
    func NewFromFunc(f HandlerFunc, name string, subs ...*Clic) *Clic
    func (c *Clic) Flag(val any, names, usage string) *flagset.Flag
    func (c *Clic) FlagRecursive(val any, names, usage string) *flagset.Flag
    func (c *Clic) HandleResolvedCmd(ctx context.Context) error
    func (c *Clic) Operand(val any, req bool, name, desc string) *operandset.Operand
    func (c *Clic) Parse(args []string) error
    func (c *Clic) SetUsageTemplating(tmplCfg *TmplConfig)
    func (c *Clic) Usage() string
type Error
    func NewError(err error, c *Clic) *Error
    func (e *Error) Clic() *Clic
    func (e *Error) Error() string
    func (e *Error) Unwrap() error
type Handler
type HandlerFunc
    func (f HandlerFunc) HandleCommand(ctx context.Context) error
type Links
    func (l Links) ParentCmd() *Clic
    func (l Links) ResolvedCmd() *Clic
    func (l Links) SubCmds() []*Clic
type ParseError
type SubRequiredError
type TmplConfig
    func NewDefaultTmplConfig() *TmplConfig
type TmplData
```

### Setup

```go
package clic_test

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/daved/clic"
)

var args = []string{"myapp", "print", "--info=flagval", "arrrg"}

func Example_noStructure() {
    var (
        info  = "default"
        value = "unset"
        out   = os.Stdout // emulate an interesting dependency
    )

    // Associate HandlerFunc with command name "print"
    print := clic.NewFromFunc(func(ctx context.Context) error {
        fmt.Fprintf(out, "info flag = %s\n", info)
        fmt.Fprintf(out, "value arg = %v\n", value)
        return nil
    }, "print")

    // Associate "print" flag and operand variables with relevant names
    print.Flag(&info, "i|info", "Set additional info.")
    print.Operand(&value, true, "first_opnd", "Value to be printed.")

    // Associate HandlerFunc with application name, adding "print" as a subcommand
    root := clic.NewFromFunc(func(ctx context.Context) error {
        fmt.Fprintln(out, "ouch, hit root")
        return nil
    }, "myapp", print)

    // Parse the cli command as `myapp print --info=flagval arrrg`
    if err := root.Parse(args[1:]); err != nil {
        log.Fatalln(err)
    }

    // Run the handler that Parse resolved to
    if err := root.HandleResolvedCmd(context.Background()); err != nil {
        log.Fatalln(err)
    }

    // Output:
    // info flag = flagval
    // value arg = arrrg
}
```

## More Info

### CLI Argument Types

There are three kinds of command line arguments that clic helps to manage: Commands/Subcommands,
Flags (plus related flag values), and Operands. Commands/subcommands each optionally have their own
flags and operands. If an argument of a command does not match a subcommand, and is not a flag arg
(i.e. it does not start with a hyphen and is not a flag value), then it will be parsed as an operand
if any operands have been defined.

Argument kinds and their placements:

```go
command --flag=flag-value subcommand -f flag-value operand_a operand_b
```

### Custom Templating

Custom templates and template behaviors (i.e. template function maps) can be set. Custom data can be
attached to instances of Clic, FlagSet, Flag, OperandSet, and Operand using their Meta fields for
access from custom templates.

### Easily Maturable

[Package docs](https://pkg.go.dev/github.com/daved/clic) contain suggestions for three stages of
application growth.
