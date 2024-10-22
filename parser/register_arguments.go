package parser

import (
	"fmt"

	"github.com/p0dalirius/goopts/argumentgroup"
	"github.com/p0dalirius/goopts/arguments"
)

// Register adds a new argument to the ArgumentsParser's default group.
// It ensures that the short and long names for the argument are unique
// and prevents duplicate registrations.
//
// Parameters:
// - arg: An instance of the arguments.Argument interface representing the argument to be registered.
//
// Returns:
// - An error if the argument's short or long name conflicts with an existing argument, otherwise nil.
func (ap *ArgumentsParser) Register(arg arguments.Argument) error {
	if ap.DefaultGroup == nil {
		ap.DefaultGroup = &argumentgroup.ArgumentGroup{}
	}

	// Initiate the maps
	if ap.DefaultGroup.ShortNameToArgument == nil {
		ap.DefaultGroup.ShortNameToArgument = make(map[string]arguments.Argument)
	}
	if ap.DefaultGroup.LongNameToArgument == nil {
		ap.DefaultGroup.LongNameToArgument = make(map[string]arguments.Argument)
	}

	if _, exists := ap.DefaultGroup.ShortNameToArgument[arg.GetShortName()]; exists {
		return fmt.Errorf("argument with short name %s already exists", arg.GetShortName())
	}
	ap.DefaultGroup.ShortNameToArgument[arg.GetShortName()] = arg

	if _, exists := ap.DefaultGroup.LongNameToArgument[arg.GetLongName()]; exists {
		return fmt.Errorf("argument with long name %s already exists", arg.GetLongName())
	}
	ap.DefaultGroup.LongNameToArgument[arg.GetLongName()] = arg

	ap.DefaultGroup.Arguments = append(ap.DefaultGroup.Arguments, arg)
	return nil
}

// NewBoolArgument initializes a new BoolArgument and registers it with the ArgumentsParser.
// It sets up the argument with the provided short and long names, default value, and help message.
//
// Parameters:
// - ptr: A pointer to the boolean variable where the argument's value will be stored.
// - shortName: The short flag (e.g., "-b") used to specify the boolean argument. It can be empty if no short flag is defined.
// - longName: The long flag (e.g., "--bool") used to specify the boolean argument. It can be empty if no long flag is defined.
// - defaultValue: The boolean value to be used if the argument is not provided by the user.
// - help: A description of what this argument represents, displayed in help/usage information.
//
// Returns:
// - An error if the argument registration fails, otherwise nil.
func (ap *ArgumentsParser) NewBoolArgument(ptr *bool, shortName, longName string, defaultValue bool, help string) error {
	arg := &arguments.BoolArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, help)
	err := ap.Register(arg)
	return err
}

// NewStringArgument initializes a new StringArgument and registers it with the ArgumentsParser.
// It sets up the argument with the provided short and long names, default value, requirement status,
// and help message.
//
// Parameters:
// - ptr: A pointer to the string variable where the argument's value will be stored.
// - shortName: The short flag (e.g., "-s") used to specify the string argument. It can be empty if no short flag is defined.
// - longName: The long flag (e.g., "--string") used to specify the string argument. It can be empty if no long flag is defined.
// - defaultValue: The string value to be used if the argument is not provided by the user.
// - required: Indicates whether the argument must be specified by the user.
// - help: A description of what this argument represents, displayed in help/usage information.
//
// Returns:
// - An error if the argument registration fails, otherwise nil.
func (ap *ArgumentsParser) NewStringArgument(ptr *string, shortName, longName string, defaultValue string, required bool, help string) error {
	arg := &arguments.StringArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, required, help)
	err := ap.Register(arg)
	return err
}

// NewIntArgument initializes a new IntArgument and registers it with the ArgumentsParser.
// It sets up the argument with the provided short and long names, default value, requirement status,
// and help message.
//
// Parameters:
// - ptr: A pointer to the integer variable where the argument's value will be stored.
// - shortName: The short flag (e.g., "-n") used to specify the integer argument. It can be empty if no short flag is defined.
// - longName: The long flag (e.g., "--number") used to specify the integer argument. It can be empty if no long flag is defined.
// - defaultValue: The integer value to be used if the argument is not provided by the user.
// - required: Indicates whether the argument must be specified by the user.
// - help: A description of what this argument represents, displayed in help/usage information.
//
// Returns:
// - An error if the argument registration fails, otherwise nil.
func (ap *ArgumentsParser) NewIntArgument(ptr *int, shortName, longName string, defaultValue int, required bool, help string) error {
	arg := &arguments.IntArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, required, help)
	err := ap.Register(arg)
	return err
}

// NewIntRangeArgument initializes a new IntRangeArgument and registers it with the ArgumentsParser.
// It sets up the argument with the provided short and long names, default value, range boundaries,
// requirement status, and help message.
//
// Parameters:
// - ptr: A pointer to the integer variable where the argument's value will be stored.
// - shortName: The short flag (e.g., "-r") used to specify the argument. It can be empty if no short flag is defined.
// - longName: The long flag (e.g., "--range") used to specify the argument. It can be empty if no long flag is defined.
// - defaultValue: The integer to be used if the argument is not provided by the user.
// - rangeStart: The minimum allowable value for the integer argument.
// - rangeStop: The maximum allowable value for the integer argument.
// - required: Indicates whether the argument must be specified by the user.
// - help: A description of what this argument represents, displayed in help/usage information.
//
// Returns:
// - An error if the argument registration fails, otherwise nil.
func (ap *ArgumentsParser) NewIntRangeArgument(ptr *int, shortName, longName string, defaultValue int, rangeStart int, rangeStop int, required bool, help string) error {
	arg := &arguments.IntRangeArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, rangeStart, rangeStop, required, help)
	err := ap.Register(arg)
	return err
}

// NewTcpPortArgument initializes a new TcpPortArgument and registers it with the ArgumentsParser.
// It sets up the argument with the provided short and long names, default value, requirement status,
// and help message.
//
// Parameters:
// - ptr: A pointer to the integer variable where the argument's value will be stored.
// - shortName: The short flag (e.g., "-p") used to specify the TCP port argument. It can be empty if no short flag is defined.
// - longName: The long flag (e.g., "--port") used to specify the TCP port argument. It can be empty if no long flag is defined.
// - defaultValue: The TCP port number to be used if the argument is not provided by the user.
// - required: Indicates whether the argument must be specified by the user.
// - help: A description of what this argument represents, displayed in help/usage information.
//
// Returns:
// - An error if the argument registration fails, otherwise nil.
func (ap *ArgumentsParser) NewTcpPortArgument(ptr *int, shortName, longName string, defaultValue int, required bool, help string) error {
	arg := &arguments.TcpPortArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, required, help)
	err := ap.Register(arg)
	return err
}
