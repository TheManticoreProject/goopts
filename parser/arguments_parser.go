package parser

import (
	"github.com/TheManticoreProject/goopts/argumentgroup"
	"github.com/TheManticoreProject/goopts/arguments"
	"github.com/TheManticoreProject/goopts/positionals"
)

// ArgumentsParser is responsible for parsing command-line arguments for a program.
// It handles both positional and named arguments, supports argument groups, and provides
// usage information. The parser maintains a mapping of argument names to their corresponding
// values and manages required arguments.
type ArgumentsParser struct {
	// Banner is a string that is printed before running the argument parsing.
	// This is typically used to display the program name or purpose.
	Banner string

	// Options holds various configuration options for the ArgumentsParser.
	Options ArgumentsParserOptions

	// ParsingState holds the state of the parsing process.
	ParsingState ParsingState

	// PositionalArguments is a slice of pointers to positional arguments
	// that the parser will manage.
	PositionalArguments []positionals.PositionalArgument

	// Groups is a map of named subgroups to organize related arguments
	// into logical categories for better structure and readability.
	Groups map[string]*argumentgroup.ArgumentGroup

	// SubParsers holds the subparsers for handling subcommands within the main parser.
	SubParsers SubParsers

	// shortNameToArgument is a map that associates short flag names (e.g., "-v")
	// with their corresponding argument structures.
	shortNameToArgument map[string]arguments.Argument

	// longNameToArgument is a map that associates long flag names (e.g., "--verbose")
	// with their corresponding argument structures.
	longNameToArgument map[string]arguments.Argument

	// requiredArguments is a slice of pointers to required arguments that must be provided
	// by the user for the program to run successfully.
	requiredArguments []arguments.Argument

	// allArguments is a slice containing all arguments (both positional and named) that
	// the parser manages.
	allArguments []arguments.Argument
}
