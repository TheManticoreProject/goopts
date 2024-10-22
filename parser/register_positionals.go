package parser

import (
	"fmt"
	"strings"

	"github.com/p0dalirius/goopts/positionals"
)

// Register registers a new argument with the argument group if it does not already exist.
//
// Parameters:
// - arg: The argument to be registered.
//
// The function checks if the argument's short name and long name are already present in the
// ShortNameToArgument and LongNameToArgument maps. If not, it adds the argument to these maps
// and appends it to the Arguments slice.
func (ap *ArgumentsParser) RegisterPositional(newPositionalArgument positionals.PositionalArgument) error {
	// fmt.Printf("[debug] ArgumentsParser.RegisterPositional(arg positionals.PositionalArgument)\n")

	if ap.PositionalArguments == nil {
		ap.PositionalArguments = make([]positionals.PositionalArgument, 0)
	}

	for _, existingPositionalArgument := range ap.PositionalArguments {
		if strings.EqualFold(existingPositionalArgument.GetName(), newPositionalArgument.GetName()) {
			return fmt.Errorf("positional argument with name %s already exists", newPositionalArgument.GetName())
		}
	}

	ap.PositionalArguments = append(ap.PositionalArguments, newPositionalArgument)

	return nil
}

// NewStringArgument registers a new string argument with the argument parser.
//
// Parameters:
// - ptr: A pointer to the string variable where the argument value will be stored.
// - shortName: The short name (single character) of the argument, prefixed with a dash (e.g., "-u").
// - longName: The long name of the argument, prefixed with two dashes (e.g., "--username").
// - defaultValue: The default value of the argument if it is not provided by the user.
// - help: A description of the argument, which will be displayed in the help message.
//
// The function creates a new StringArgument with the provided parameters and adds it to the argument group.
func (ap *ArgumentsParser) NewStringPositionalArgument(ptr *string, name string, help string) error {
	arg := &positionals.StringPositionalArgument{}
	arg.Init(ptr, name, help)
	err := ap.RegisterPositional(arg)
	return err
}

// NewIntArgument registers a new int argument with the argument parser.
//
// Parameters:
// - ptr: A pointer to the int variable where the argument value will be stored.
// - shortName: The short name (single character) of the argument, prefixed with a dash (e.g., "-u").
// - longName: The long name of the argument, prefixed with two dashes (e.g., "--username").
// - defaultValue: The default value of the argument if it is not provided by the user.
// - help: A description of the argument, which will be displayed in the help message.
//
// The function creates a new StringArgument with the provided parameters and adds it to the argument group.
func (ap *ArgumentsParser) NewIntPositionalArgument(ptr *int, name string, help string) error {
	arg := &positionals.IntPositionalArgument{}
	arg.Init(ptr, name, help)
	err := ap.RegisterPositional(arg)
	return err
}
