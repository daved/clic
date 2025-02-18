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
    func (c *Clic) Handle(ctx context.Context) error
    func (c *Clic) Operand(val any, req bool, name, desc string) *operandset.Operand
    func (c *Clic) Parse(args []string) error
    func (c *Clic) Recursively(fn func(*Clic))
    func (c *Clic) Usage() string
// see package docs for more
```

### Setup

```go
func main() {
    // error handling omitted to keep example focused

    var (
        info  = "default"
        value = "unset"
    )

    // Associate HandlerFunc with command name
    root := clic.NewFromFunc(printFunc(&info, &value), "myapp")

    // Associate flag and operand variables with relevant names
    root.Flag(&info, "i|info", "Set additional info.")
    root.Operand(&value, true, "first_operand", "Value to be printed.")

    // Parse the cli command as `myapp --info=flagval arrrg`
    cmd, _ := c.Parse([]string{"--info=flagval", "arrrg"})

    // Run the handler that Parse resolved to
    _ = cmd.Handle(context.Background())
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

### Default Templating

`cmd.Usage()` value from the usage example above:

```txt
Usage:

  myapp [FLAGS] <first_operand>

Flags for myapp:

    -i, --info  =STRING    default: default
        Set additional info.
```

### Custom Templating

Custom templates and template behaviors (i.e. template function maps) can be used with the clic/tmpl
package for convenience. Custom data can be attached to instances of Clic, FlagSet, Flag,
OperandSet, and Operand using their Meta fields for access from custom templates. The contents of
NewUsageTmpl can be used as a reference.

### Maturable Architecture

[Package docs](https://pkg.go.dev/github.com/daved/clic) contain suggestions for three stages of
application growth.

### Dependencies

- [flagset](https://github.com/daved/flagset)
- [operandset](https://github.com/daved/operandset)
- [vtype](https://github.com/daved/vtype)
