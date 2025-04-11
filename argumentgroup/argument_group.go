package argumentgroup

import (
	"fmt"
	"strings"

	"github.com/TheManticoreProject/goopts/arguments"
)

const (
	// ARGUMENT_GROUP_TYPE_NORMAL defines a basic group of arguments
	// where each argument is treated independently. There are no specific
	// rules governing which arguments must or must not be set.
	// This is the default behavior for argument groups.
	ARGUMENT_GROUP_TYPE_NORMAL = 0

	// ARGUMENT_GROUP_TYPE_REQUIRED_MUTUALLY_EXCLUSIVE creates a group
	// where exactly one argument must be set. If more than one or none
	// of the arguments in this group are specified, the parser will throw an error.
	// This is useful for scenarios where a user must choose one of several options.
	ARGUMENT_GROUP_TYPE_REQUIRED_MUTUALLY_EXCLUSIVE = 1

	// ARGUMENT_GROUP_TYPE_NOT_REQUIRED_MUTUALLY_EXCLUSIVE defines a group
	// where at most one argument can be set, but it is not mandatory to set any.
	// If multiple arguments from this group are provided, the parser will throw an error.
	// It is designed for optional settings where only one should be active if chosen.
	ARGUMENT_GROUP_TYPE_NOT_REQUIRED_MUTUALLY_EXCLUSIVE = 2

	// ARGUMENT_GROUP_TYPE_DEPENDENT ensures that if one argument in the group
	// is set, all other arguments in the group must also be set. If any are missing,
	// the parser will throw an error. This is useful for scenarios where arguments
	// are interdependent and must be specified together for correct operation.
	ARGUMENT_GROUP_TYPE_DEPENDENT = 3
)

// ArgumentGroup represents a group of arguments with associated metadata.
// It is used to categorize and manage related command-line arguments.
type ArgumentGroup struct {
	// Name is the identifier for the argument group, often used for display purposes.
	Name string

	// Arguments is a list of Argument pointers that belong to this group.
	Arguments []arguments.Argument

	// ShortNameToArgument maps short argument names (e.g., "-f") to their corresponding Argument.
	ShortNameToArgument map[string]arguments.Argument

	// LongNameToArgument maps long argument names (e.g., "--file") to their corresponding Argument.
	LongNameToArgument map[string]arguments.Argument

	// Type is an integer that can represent the type or category of the argument group.
	// This can be used for conditional handling or to differentiate between groups.
	Type int
}

// Register registers a new argument with the argument group if it does not already exist.
//
// Parameters:
// - arg: The argument to be registered.
//
// The function checks if the argument's short name and long name are already present in the
// ShortNameToArgument and LongNameToArgument maps. If not, it adds the argument to these maps
// and appends it to the Arguments slice.
func (ag *ArgumentGroup) Register(arg arguments.Argument) error {
	if ag.ShortNameToArgument == nil {
		ag.ShortNameToArgument = make(map[string]arguments.Argument)
	}
	if ag.LongNameToArgument == nil {
		ag.LongNameToArgument = make(map[string]arguments.Argument)
	}

	if len(arg.GetShortName()) != 0 {
		if _, exists := ag.ShortNameToArgument[arg.GetShortName()]; exists {
			return fmt.Errorf("argument with short name %s already exists", arg.GetShortName())
		}
		ag.ShortNameToArgument[arg.GetShortName()] = arg
	}

	if len(arg.GetLongName()) != 0 {
		if _, exists := ag.LongNameToArgument[arg.GetLongName()]; exists {
			return fmt.Errorf("argument with long name %s already exists", arg.GetLongName())
		}
		ag.LongNameToArgument[arg.GetLongName()] = arg
	}

	ag.LongNameToArgument[arg.GetLongName()] = arg

	ag.Arguments = append(ag.Arguments, arg)
	return nil
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
// with a dash. It then checks if the argument is in the `LongNameToArgument` map
// for long arguments or `ShortNameToArgument` map for short arguments, returning
// true if found and false if not.
func (ag *ArgumentGroup) ArgumentIsPresent(argumentName string) bool {
	if len(argumentName) < 2 || argumentName[0] != '-' {
		return false
	}

	if argumentName[1] == '-' {
		// Long argument flag (e.g., --example)
		if argument, exists := ag.LongNameToArgument[argumentName]; exists {
			return argument.IsPresent()
		}
	} else {
		// Short argument flag (e.g., -e)
		if argument, exists := ag.ShortNameToArgument[argumentName]; exists {
			return argument.IsPresent()
		}
	}

	// Argument not found
	return false
}

// PrintArgumentTree prints the argument tree for the argument group.
//
// Parameters:
// - indent: The indentation level for the printed output.
//
// The function prints the name of the argument group and its arguments in a tree-like structure,
// with each level of indentation represented by "  │ ". The output includes the group name and
// the names of the arguments within the group.
func (ag *ArgumentGroup) PrintArgumentTree(indent int) {
	indentPrompt := strings.Repeat("  │ ", indent)
	//
	fmt.Printf("%s  ├─ <Group name=\"%s\">\n", indentPrompt, ag.Name)

	fmt.Printf("%s  │   ├─ Name: \"%s\"\n", indentPrompt, ag.Name)

	fmt.Printf("%s  │   ├─ Arguments (%d): \n", indentPrompt, len(ag.Arguments))
	for _, argument := range ag.Arguments {
		argtype := ""
		if _, ok := argument.(*arguments.BoolArgument); ok {
			argtype = "bool"
		} else if _, ok := argument.(*arguments.StringArgument); ok {
			argtype = "string"
		} else if _, ok := argument.(*arguments.IntArgument); ok {
			argtype = "int"
		}
		fmt.Printf("%s  │   │   ├─ (\"%s\",\"%s\") [%s] \"%s\"\n", indentPrompt, argument.GetShortName(), argument.GetLongName(), argtype, argument.GetHelp())
	}
	fmt.Printf("%s  │   │   └──\n", indentPrompt)
	fmt.Printf("%s  │   └──\n", indentPrompt)
	fmt.Printf("%s  └──\n", indentPrompt)
}
