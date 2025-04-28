package parser

import (
	"fmt"

	"github.com/TheManticoreProject/goopts/argumentgroup"
	"github.com/TheManticoreProject/goopts/arguments"
)

func NewParser(banner string) *ArgumentsParser {
	return &ArgumentsParser{
		Banner: banner,
		Groups: make(map[string]*argumentgroup.ArgumentGroup),
	}
}

// ArgumentIsPresent checks if a given argument is present in the parsed arguments.
// It supports both short (e.g., -e) and long (e.g., --example) argument names.
//
// Parameters:
//   - argumentName: The name of the argument to check. It should start with a dash ('-')
//     for short arguments or two dashes ('--') for long arguments.
//
// Returns:
// - bool: true if the argument is present; false otherwise.
//
// The function first verifies that the argument name has a valid length and starts
// with a dash. It then checks if the argument is in the `longNameToArgument` map
// for long arguments or `shortNameToArgument` map for short arguments, returning
// true if found and false if not.
func (ap *ArgumentsParser) ArgumentIsPresent(argumentName string) bool {
	if len(argumentName) < 2 || argumentName[0] != '-' {
		return false
	}

	if argumentName[1] == '-' {
		// Long argument flag (e.g., --example)
		if argument, exists := ap.ParsingState.ParsedArguments.LongNameToArgument[argumentName]; exists {
			return (*argument).IsPresent()
		}
	} else {
		// Short argument flag (e.g., -e)
		if argument, exists := ap.ParsingState.ParsedArguments.ShortNameToArgument[argumentName]; exists {
			return (*argument).IsPresent()
		}
	}

	// Argument not found
	return false
}

// Get retrieves the value of the argument specified by its short or long name.
//
// Parameters:
//   - name: The name of the argument to retrieve. This can be either the short name (single character)
//     or the long name (string) of the argument.
//
// Returns:
// - The value of the argument as an interface{} if the argument is found.
// - An error if the argument maps are not initialized or if the argument with the specified name is not found.
//
// The function first checks if the argument maps (shortNameToArgument and longNameToArgument) are initialized.
// If they are not, it returns an error indicating that the argument maps are not initialized.
// It then looks up the argument by its short name in the shortNameToArgument map.
// If the argument is not found, it looks up the argument by its long name in the longNameToArgument map.
// If the argument is still not found, it returns an error indicating that the argument with the specified name is not found.
// If the argument is found, it returns the value of the argument.
func (ap *ArgumentsParser) Get(argumentFlag string) (interface{}, error) {
	if ap.shortNameToArgument == nil || ap.longNameToArgument == nil {
		return nil, fmt.Errorf("argument maps are not initialized")
	}

	var arg arguments.Argument
	var exists bool

	if arg, exists = ap.shortNameToArgument[argumentFlag]; !exists {
		if arg, exists = ap.longNameToArgument[argumentFlag]; !exists {
			return nil, fmt.Errorf("argument '%s' not found", argumentFlag)
		}
	}

	return arg.GetValue(), nil
}

// PrintArgumentTree prints the argument tree for the ArgumentsParser.
//
// The function prints the banner, the list of arguments, and the subgroups in a tree-like structure.
// Each level of indentation is represented by "  │ ". The output includes the banner, the names of the arguments,
// and the names of the subgroups within the ArgumentsParser.
func (ap *ArgumentsParser) PrintArgumentTree() {
	indent := 0
	//
	fmt.Printf("<ArgumentsParser>\n")

	fmt.Printf("  ├─ Banner: \"%s\"\n", ap.Banner)

	fmt.Printf("  ├─ Arguments (%d): \n", len(ap.Groups[""].Arguments))
	for _, argument := range ap.Groups[""].Arguments {
		argtype := ""
		if _, ok := argument.(*arguments.BoolArgument); ok {
			argtype = "bool"
		} else if _, ok := argument.(*arguments.StringArgument); ok {
			argtype = "string"
		} else if _, ok := argument.(*arguments.IntArgument); ok {
			argtype = "int"
		}
		fmt.Printf("  │   │   ├─ (\"%s\",\"%s\") [%s] \"%s\"\n", argument.GetShortName(), argument.GetLongName(), argtype, argument.GetHelp())
	}
	fmt.Printf("  │   │   └──\n")
	fmt.Printf("  │   └──\n")

	// Call subgroups
	fmt.Printf("  ├─ SubGroups: \n")
	for _, subgroup := range ap.Groups {
		subgroup.PrintArgumentTree(indent + 1)
	}
	fmt.Printf("  └──\n")
}
