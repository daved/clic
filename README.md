# clic [![GoDoc](https://pkg.go.dev/badge/github.com/daved/clic.svg)](https://pkg.go.dev/github.com/daved/clic)

```sh
go get github.com/daved/clic
```

Package clic provides streamlined POSIX-friendly CLI command parsing.

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

Three types of arguments are handled: commands (and subcommands), flags (and their values), and
operands. Commands and subcommands can each define their own flags and operands. Arguments that are
not subcommands or flag-related (no hyphen prefix, not a flag value), will be treated as operands.

Example argument layout:

```sh
command --flag=flag-value subcommand -f flag-value operand_a operand_b
```

### Default Templating

`cmd.Usage()` value from the example above:

```txt
Usage:

  myapp [FLAGS] <first_operand>

Flags for myapp:

    -i, --info  =STRING    default: default
        Set additional info.
```

### Custom Templating

The Tmpl type eases custom templating. Custom data can be attached to instances of Clic, FlagSet,
Flag, OperandSet, and Operand via their Meta fields for use in templates. The default template
construction function (NewUsageTmpl) can be used as a reference for custom templates.

### Designed To Be...

POSIX-friendly:

- Long and short flags should be supported in a reasonably normal manner
- Use the term "Operand" rather than the ambiguous "Arg"
- Flags should not be processed after operands
- Usage output should be familiar and effective

Simple and powerful:

- The API surface should be small and flexible
- Flags and operands should not be auto-set (e.g. help flag set explicitly)
- Advanced usage should look similar to basic usage
- Users should be trusted to understand/use the type system (e.g. interfaces)
- Users should be trusted to understand/use advanced control flow mechanisms

### Maturable Architecture

[Package docs](https://pkg.go.dev/github.com/daved/clic) contain suggestions for three stages of
application growth.

### Dependencies

- [flagset](https://github.com/daved/flagset)
- [operandset](https://github.com/daved/operandset)
- [vtype](https://github.com/daved/vtype)
