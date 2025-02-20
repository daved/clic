// Package clic provides streamlined POSIX-friendly CLI command parsing.
//
// # CLI Argument Types
//
// Three types of arguments are handled: commands (and subcommands), flags (and
// their values), and operands. Commands and subcommands can each define their
// own flags and operands. Arguments that are not subcommands or flag-related
// (no hyphen prefix, not a flag value) will be treated as operands.
//
// Example argument layout:
//
//	command --flag=flag-value subcommand -f flag-value operand_a operand_b
//
// # Custom Templating
//
// The [Tmpl] type eases custom templating. Custom data can be attached to
// instances of Clic, FlagSet, Flag, OperandSet, and Operand via their Meta
// fields for use in templates. The default template construction function
// ([NewUsageTmpl]) can be used as a reference for custom templates.
//
// # Designed To Be...
//
// POSIX-friendly:
//
//   - Long and short flags should be supported in a reasonably normal manner
//   - Use the term "Operand" rather than the ambiguous "Arg"
//   - Flags should not be processed after operands
//   - Usage output should be familiar and effective
//
// Simple and powerful:
//
//   - The API surface should be small and flexible
//   - Flags and operands should not be auto-set (e.g. help flag set explicitly)
//   - Advanced usage should look similar to basic usage
//   - Users should be trusted to understand/use the type system (e.g. interfaces)
//   - Users should be trusted to understand/use advanced control flow mechanisms
package clic
