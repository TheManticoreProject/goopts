package arguments

import (
	"fmt"

	"github.com/p0dalirius/goopts/utils"
)

// ListOfIntsArgument represents a command-line argument that expects a list of integers.
// It provides information about the argument's short and long flag names, help message,
// the default value, and whether the argument is required.
type ListOfIntsArgument struct {
	// ShortName is the short flag (e.g., "-i") used to specify the list of integers.
	// It can be empty if no short flag is defined.
	ShortName string
	// LongName is the long flag (e.g., "--ints") used to specify the list of integers.
	// It can be empty if no long flag is defined.
	LongName string
	// Help provides a description of what this argument represents.
	// This message is displayed when showing help/usage information.
	Help string
	// Value stores the actual list of integers provided by the user.
	// If no value is specified by the user, Value will hold the DefaultValue.
	Value *[]int
	// DefaultValue is the list of integers to be used if the argument is not provided by the user.
	DefaultValue []int
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
func (arg ListOfIntsArgument) GetShortName() string {
	return arg.ShortName
}

// GetLongName returns the long flag name of the argument.
// If no long flag is defined, it returns an empty string.
func (arg ListOfIntsArgument) GetLongName() string {
	return arg.LongName
}

// GetHelp returns the help message of the argument.
// This provides a description of how to use the argument.
func (arg ListOfIntsArgument) GetHelp() string {
	return arg.Help
}

// GetValue returns the current list of integers as an interface{}.
// It will return the actual value provided by the user or the default value if none was specified.
func (arg ListOfIntsArgument) GetValue() any {
	return *arg.Value
}

// GetDefaultValue returns the default list of integers as an interface{}.
// This is used when the argument is not specified by the user.
func (arg ListOfIntsArgument) GetDefaultValue() any {
	return arg.DefaultValue
}

// ResetDefaultValue resets the value of the argument to the default value.
func (arg *ListOfIntsArgument) ResetDefaultValue() {
	*(arg.Value) = append(*(arg.Value), arg.DefaultValue...)
}

// IsRequired returns whether the argument is required.
// If true, the argument must be specified when running the program.
func (arg ListOfIntsArgument) IsRequired() bool {
	return arg.Required
}

// IsPresent checks if the argument was set in the command line.
func (arg ListOfIntsArgument) IsPresent() bool {
	return arg.Present
}

// Init initializes the ListOfIntsArgument with the provided parameters.
// It sets the flag names, required status, help message, actual value, and default value.
func (arg *ListOfIntsArgument) Init(value *[]int, shortName, longName string, defaultValue []int, required bool, help string) {
	arg.LongName, arg.ShortName = utils.GenerateLongAndShortNames(longName, shortName)

	arg.Required = required

	arg.Present = false

	arg.Help = help

	arg.Value = value

	arg.DefaultValue = defaultValue
}

// Consume processes the command-line arguments and sets the value of the ListOfIntsArgument.
//
// The function iterates through the provided arguments and checks if any of them match the short or long name of the ListOfIntsArgument.
// If a match is found, it sets the value of the ListOfIntsArgument to the next argument in the list and returns the remaining arguments.
// If the argument is required and not found, it returns an error.
//
// Parameters:
//   - arguments: A slice of strings representing the command-line arguments.
//
// Returns:
// - A slice of strings representing the remaining arguments after processing the ListOfIntsArgument.
func (arg *ListOfIntsArgument) Consume(arguments []string) ([]string, error) {
	// Initiate the value for the first time
	if arg.Value == nil {
		v := make([]int, 0)
		arg.Value = &v
	}

	sizeToConsume := 2

	if len(arguments) >= sizeToConsume {
		if (arguments[0] == arg.ShortName) || (arguments[0] == arg.LongName) {
			value, err := utils.StringToInt(arguments[1])
			if err != nil {
				// Return the original arguments if parsing fails
				return arguments, fmt.Errorf("%s %s: could not parse integer: %s", arguments[0], arguments[1], err)
			}
			*arg.Value = append(*arg.Value, int(value))

			arg.Present = true

			return arguments[sizeToConsume:], nil
		}
	}

	return arguments, nil
}
