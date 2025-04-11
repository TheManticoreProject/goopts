package argumentgroup

import (
	"github.com/TheManticoreProject/goopts/arguments"
)

// NewBoolArgument registers a new boolean argument with the argument parser.
//
// Parameters:
// - ptr: A pointer to the boolean variable where the argument value will be stored.
// - shortName: The short name (single character) of the argument, prefixed with a dash (e.g., "-S").
// - longName: The long name of the argument, prefixed with two dashes (e.g., "--use-ldaps").
// - defaultValue: The default value of the argument if it is not provided by the user.
// - help: A description of the argument, which will be displayed in the help message.
//
// The function creates a new BoolArgument with the provided parameters and adds it to the argument group.
func (ag *ArgumentGroup) NewBoolArgument(ptr *bool, shortName, longName string, defaultValue bool, help string) error {
	arg := arguments.BoolArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, help)
	err := ag.Register(&arg)
	return err
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
func (ag *ArgumentGroup) NewStringArgument(ptr *string, shortName, longName string, defaultValue string, required bool, help string) error {
	arg := arguments.StringArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, required, help)
	err := ag.Register(&arg)
	return err
}

// NewIntArgument registers a new integer argument with the argument parser.
//
// Parameters:
// - ptr: A pointer to the integer variable where the argument value will be stored.
// - shortName: The short name (single character) of the argument, prefixed with a dash (e.g., "-i").
// - longName: The long name of the argument, prefixed with two dashes (e.g., "--iterations").
// - defaultValue: The default value of the argument if it is not provided by the user.
// - help: A description of the argument, which will be displayed in the help message.
//
// The function creates a new IntArgument with the provided parameters and adds it to the argument group.
func (ag *ArgumentGroup) NewIntArgument(ptr *int, shortName, longName string, defaultValue int, required bool, help string) error {
	arg := arguments.IntArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, required, help)
	err := ag.Register(&arg)
	return err
}

// NewIntRangeArgument registers a new integer range argument with the argument parser.
//
// Parameters:
// - ptr: A pointer to the integer variable where the argument value will be stored.
// - shortName: The short name (single character) of the argument, prefixed with a dash (e.g., "-i").
// - longName: The long name of the argument, prefixed with two dashes (e.g., "--iterations").
// - defaultValue: The default value of the argument if it is not provided by the user.
// - min: The minimum value of the argument.
// - max: The maximum value of the argument.
// - required: Indicates whether the argument must be specified by the user.
// - help: A description of the argument, which will be displayed in the help message.
//
// The function creates a new IntRangeArgument with the provided parameters and adds it to the argument group.
func (ag *ArgumentGroup) NewIntRangeArgument(ptr *int, shortName, longName string, defaultValue int, min int, max int, required bool, help string) error {
	arg := arguments.IntRangeArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, min, max, required, help)
	err := ag.Register(&arg)
	return err
}

// NewTcpPortArgument registers a new integer argument with the argument parser.
//
// Parameters:
// - ptr: A pointer to the integer variable where the argument value will be stored.
// - shortName: The short name (single character) of the argument, prefixed with a dash (e.g., "-i").
// - longName: The long name of the argument, prefixed with two dashes (e.g., "--iterations").
// - defaultValue: The default value of the argument if it is not provided by the user.
// - help: A description of the argument, which will be displayed in the help message.
//
// The function creates a new IntArgument with the provided parameters and adds it to the argument group.
func (ag *ArgumentGroup) NewTcpPortArgument(ptr *int, shortName, longName string, defaultValue int, required bool, help string) error {
	arg := arguments.TcpPortArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, required, help)
	err := ag.Register(&arg)
	return err
}
