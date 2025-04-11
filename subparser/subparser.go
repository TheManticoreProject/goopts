package subparser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TheManticoreProject/goopts/parser"
)

// The ArgumentsParser is designed to handle command-line argument parsing for applications.
// It supports both positional and named arguments, organizes them into groups, and provides
// detailed usage information. The parser keeps track of argument names and their values,
// ensuring that all required arguments are accounted for.
//
// Attributes:
//
//	Banner (string): A string displayed before parsing begins, often used to show the program's name or purpose.
//	PositionalArguments ([]*positionals.PositionalArgument): A list of pointers to positional arguments managed by the parser.
//	DefaultGroup (*argumentgroup.ArgumentGroup): A pointer to the default group for organizing arguments, if applicable.
//	shortNameToArgument (map[string]arguments.Argument): A map linking short flag names (e.g., "-v") to their argument structures.
//	longNameToArgument (map[string]arguments.Argument): A map linking long flag names (e.g., "--verbose") to their argument structures.
//	requiredArguments ([]arguments.Argument): A list of pointers to arguments that must be provided for successful execution.
//	allArguments ([]arguments.Argument): A list of all arguments, both positional and named, managed by the parser.
//	Groups (map[string]*argumentgroup.ArgumentGroup): A map organizing related arguments into named subgroups for clarity and structure.
type ArgumentsSubparser struct {
	Banner           string
	ShowBannerOnHelp bool
	ShowBannerOnRun  bool

	Name  string
	Value *string

	CaseInsensitive bool

	SubParsers map[string]*parser.ArgumentsParser
}

// NewSubparser creates a new subparser with the specified name.
//
// Parameters:
// - name: The name of the subparser.
// - help: The help message of the subparser.
//
// Returns:
// - A pointer to the newly created subparser.
func (asp *ArgumentsSubparser) AddSubParser(name, banner string) *parser.ArgumentsParser {
	parser_ptr := &parser.ArgumentsParser{
		Banner:           banner,
		ShowBannerOnHelp: false,
		ShowBannerOnRun:  false,
	}

	if asp.SubParsers == nil {
		asp.SubParsers = make(map[string]*parser.ArgumentsParser)
	}

	if asp.CaseInsensitive {
		asp.SubParsers[strings.ToLower(name)] = parser_ptr
	} else {
		asp.SubParsers[name] = parser_ptr
	}

	return parser_ptr
}

// Parse parses the arguments and returns the parsed arguments.
//
// Returns:
// - A map of parsed arguments.
func (asp *ArgumentsSubparser) Parse() {
	if len(asp.Banner) != 0 {
		fmt.Printf("%s\n\n", asp.Banner)
	}

	if len(os.Args) < 2 {
		asp.Usage()
		os.Exit(1)
	} else {
		// Get the subparser name from the first argument
		subparser_name := os.Args[1]
		if asp.CaseInsensitive {
			subparser_name = strings.ToLower(os.Args[1])
		}

		// Check if the user wants to see the help message
		if subparser_name == "help" || subparser_name == "--help" || subparser_name == "-h" {
			asp.Usage()
			os.Exit(0)
		}

		// Consume the program name and the subparser name
		if subparser, exists := asp.SubParsers[subparser_name]; exists {
			if asp.Value != nil {
				*asp.Value = subparser_name
			}
			subparser.ParseFrom(2)
		} else {
			asp.Usage()
			if len(asp.Name) != 0 {
				fmt.Printf("[!] No %s with name '%s' found.\n", asp.Name, subparser_name)
			} else {
				fmt.Printf("[!] No subparser with name '%s' found.\n", subparser_name)
			}
			os.Exit(1)
		}
	}
}

// Usage prints the usage information for the subparser.
//
// Returns:
// - A map of parsed arguments.
func (asp *ArgumentsSubparser) Usage() {
	if len(asp.Name) != 0 {
		fmt.Printf("Usage: %s <%s> [options]\n", filepath.Base(os.Args[0]), asp.Name)
	} else {
		fmt.Printf("Usage: %s <subparser> [options]\n", filepath.Base(os.Args[0]))
	}

	if len(asp.SubParsers) != 0 {
		fmt.Printf("\n")

		// Compute format string for subparser names
		max_len_subparser_string := 0
		for name := range asp.SubParsers {
			max_len_subparser_string = max(max_len_subparser_string, len(name))
		}
		fmtString := fmt.Sprintf("  %%-%ds %%s\n", max_len_subparser_string+2)

		// Sort subparser names
		subparserNames := make([]string, 0, len(asp.SubParsers))
		for name := range asp.SubParsers {
			subparserNames = append(subparserNames, name)
		}
		sort.Strings(subparserNames)

		// Print subparser names and banners
		for _, name := range subparserNames {
			subparser := asp.SubParsers[name]
			fmt.Printf(fmtString, name, subparser.Banner)
		}
	}
	fmt.Printf("\n")
}
