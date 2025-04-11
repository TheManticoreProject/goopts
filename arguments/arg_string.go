package arguments

import (
	"fmt"

	"github.com/TheManticoreProject/goopts/utils"
)

// StringArgument represents a command-line argument that expects a string value.
// It contains information about the argument's short and long flag names, help message,
// the default value, and whether the argument is required.
type StringArgument struct {
	// ShortName is the short flag (e.g., "-n") used to specify the string value.
	// It can be empty if no short flag is defined.
	ShortName string
	// LongName is the long flag (e.g., "--name") used to specify the string value.
	// It can be empty if no long flag is defined.
	LongName string
	// Help provides a description of what this argument represents.
	// This message is displayed when showing help/usage information.
	Help string
	// Value stores the actual string value provided by the user.
	// If no value is specified by the user, Value will hold the DefaultValue.
	Value *string
	// DefaultValue is the string to be used if the argument is not provided by the user.
	DefaultValue string
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
func (arg StringArgument) GetShortName() string {
	return arg.ShortName
}

// GetLongName returns the long flag name of the argument.
// If no long flag is defined, it returns an empty string.
func (arg StringArgument) GetLongName() string {
	return arg.LongName
}

// GetHelp returns the help message of the argument.
// If the argument is optional, it appends the default value to the message.
func (arg StringArgument) GetHelp() string {
	if !arg.IsRequired() {
		return fmt.Sprintf("%s (default: \"%v\")", arg.Help, arg.GetDefaultValue())
	} else {
		return arg.Help
	}
}

// GetValue returns the current value of the argument as an interface{}.
// It will return the actual value provided by the user or the default value if none was specified.
func (arg StringArgument) GetValue() any {
	return *arg.Value
}

// SetValue sets the value of the StringArgument.
// This is the string provided by the user or set by default.
func (arg *StringArgument) SetValue(value any) {
	*(arg.Value) = value.(string)
}

// GetDefaultValue returns the default value of the argument as an interface{}.
// This is used when the argument is not specified by the user.
func (arg StringArgument) GetDefaultValue() any {
	return arg.DefaultValue
}

// ResetDefaultValue resets the value of the argument to the default value.
func (arg *StringArgument) ResetDefaultValue() {
	*(arg.Value) = arg.DefaultValue
}

// IsRequired returns whether the argument is required.
// If true, the argument must be specified when running the program.
func (arg StringArgument) IsRequired() bool {
	return arg.Required
}

// IsPresent checks if the argument was set in the command line.
func (arg StringArgument) IsPresent() bool {
	return arg.Present
}

// Init initializes the StringArgument with the provided values.
//
// The function sets the short name, long name, help message, value, and default value for the StringArgument.
// It ensures that the short name and long name are properly formatted with leading dashes.
//
// Parameters:
//   - value: A pointer to a string where the value of the argument will be stored.
//   - shortName: The short name of the argument (single character). If empty, it will be set to an empty string.
//   - longName: The long name of the argument (string). If empty, it will be set to an empty string.
//   - defaultValue: The default value of the argument.
//   - help: The help message describing the argument.
//
// The function uses the utils.StripLeftDashes function to remove any leading dashes from the short name and long name
// before adding a single dash for the short name and a double dash for the long name.
func (arg *StringArgument) Init(value *string, shortName, longName string, defaultValue string, required bool, help string) {
	arg.LongName, arg.ShortName = utils.GenerateLongAndShortNames(longName, shortName)

	arg.Required = required

	arg.Present = false

	arg.Help = help

	arg.Value = value

	arg.DefaultValue = defaultValue
}

// Consume processes the command-line arguments and sets the value of the StringArgument.
//
// The function iterates through the provided arguments and checks if any of them match the short or long name of the StringArgument.
// If a match is found, it sets the value of the StringArgument to the next argument in the list and returns the remaining arguments.
// If the argument is required and not found, it returns an error.
//
// Parameters:
//   - arguments: A slice of strings representing the command-line arguments.
//
// Returns:
// - A slice of strings representing the remaining arguments after processing the StringArgument.
func (arg *StringArgument) Consume(arguments []string) ([]string, error) {
	sizeToConsume := 2

	if len(arguments) >= sizeToConsume {
		if (arguments[0] == arg.ShortName) || (arguments[0] == arg.LongName) {
			*arg.Value = arguments[1]

			arg.Present = true

			return arguments[sizeToConsume:], nil
		}
	}

	return arguments, nil
}
