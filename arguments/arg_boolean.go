package arguments

import (
	"github.com/TheManticoreProject/goopts/utils"
)

// BoolArgument represents a command-line argument that expects a boolean value.
// It provides information about the argument's short and long flag names, help message,
// the default value, and whether the argument is required.
type BoolArgument struct {
	// ShortName is the short flag (e.g., "-b") used to specify the boolean argument.
	// It can be empty if no short flag is defined.
	ShortName string
	// LongName is the long flag (e.g., "--boolean") used to specify the boolean argument.
	// It can be empty if no long flag is defined.
	LongName string
	// Help provides a description of what this argument represents.
	// This message is displayed when showing help/usage information.
	Help string
	// Value stores the actual boolean value provided by the user.
	// If no value is specified by the user, Value will hold the DefaultValue.
	Value *bool
	// DefaultValue is the boolean to be used if the argument is not provided by the user.
	DefaultValue bool
	// Required indicates whether this argument must be specified by the user.
	// If true, the argument must be included when running the program.
	Required bool
	// Present indicates whether this argument was set by the user during execution.
	// This can be used to differentiate between arguments that were provided and those that were not,
	// allowing for different handling of default values or other logic in the program.
	Present bool
}

// GetShortName returns the short flag name of the argument.
// If no short flag is defined, it returns an empty string.
func (arg BoolArgument) GetShortName() string {
	return arg.ShortName
}

// GetLongName returns the long flag name of the argument.
// If no long flag is defined, it returns an empty string.
func (arg BoolArgument) GetLongName() string {
	return arg.LongName
}

// GetHelp returns the help message of the argument.
// This provides a description of how to use the argument.
func (arg BoolArgument) GetHelp() string {
	return arg.Help
}

// GetValue returns the current boolean value as an interface{}.
// It will return the actual value provided by the user or the default value if none was specified.
func (arg BoolArgument) GetValue() any {
	return *arg.Value
}

// SetValue sets the value of the BoolArgument.
// This is the boolean provided by the user or set by default.
func (arg *BoolArgument) SetValue(value any) {
	*(arg.Value) = value.(bool)
}

// GetDefaultValue returns the default boolean value as an interface{}.
// This is used when the argument is not specified by the user.
func (arg BoolArgument) GetDefaultValue() any {
	return arg.DefaultValue
}

// ResetDefaultValue resets the value of the argument to the default value.
func (arg *BoolArgument) ResetDefaultValue() {
	*(arg.Value) = arg.DefaultValue
}

// IsRequired returns whether the argument is required.
// If true, the argument must be specified when running the program.
func (arg BoolArgument) IsRequired() bool {
	return arg.Required
}

// IsPresent checks if the argument was set in the command line.
func (arg BoolArgument) IsPresent() bool {
	return arg.Present
}

// Init initializes the BoolArgument with the provided parameters.
// It sets the flag names, help message, actual value, and default value.
func (arg *BoolArgument) Init(value *bool, shortName, longName string, defaultValue bool, help string) {
	arg.LongName, arg.ShortName = utils.GenerateLongAndShortNames(longName, shortName)

	arg.Help = help

	arg.Present = false

	arg.Value = value

	arg.DefaultValue = defaultValue
}

// Consume processes the command-line arguments and sets the value of the BoolArgument.
//
// The function iterates through the provided arguments and checks if any of them match the short or long name of the BoolArgument.
// If a match is found, it sets the value of the BoolArgument to true and returns the remaining arguments.
// If the argument is required and not found, it returns an error.
//
// Parameters:
//   - arguments: A slice of strings representing the command-line arguments.
//
// Returns:
// - A slice of strings representing the remaining arguments after processing the BoolArgument.
func (arg *BoolArgument) Consume(arguments []string) ([]string, error) {
	sizeToConsume := 1

	if len(arguments) >= sizeToConsume {
		if (arguments[0] == arg.ShortName) || (arguments[0] == arg.LongName) {
			*(arg.Value) = !arg.DefaultValue

			arg.Present = true

			return arguments[sizeToConsume:], nil
		}
	}

	return arguments, nil
}
