package arguments

import "github.com/p0dalirius/goopts/utils"

// ListOfStringsArgument represents a command-line argument that expects a list of strings.
// It provides information about the argument's short and long flag names, help message,
// the default value, and whether the argument is required.
type ListOfStringsArgument struct {
	// ShortName is the short flag (e.g., "-s") used to specify the list of strings.
	// It can be empty if no short flag is defined.
	ShortName string
	// LongName is the long flag (e.g., "--strings") used to specify the list of strings.
	// It can be empty if no long flag is defined.
	LongName string
	// Help provides a description of what this argument represents.
	// This message is displayed when showing help/usage information.
	Help string
	// Value stores the actual list of strings provided by the user.
	// If no value is specified by the user, Value will hold the DefaultValue.
	Value *[]string
	// DefaultValue is the list of strings to be used if the argument is not provided by the user.
	DefaultValue []string
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
func (arg ListOfStringsArgument) GetShortName() string {
	return arg.ShortName
}

// GetLongName returns the long flag name of the argument.
// If no long flag is defined, it returns an empty string.
func (arg ListOfStringsArgument) GetLongName() string {
	return arg.LongName
}

// GetHelp returns the help message of the argument.
// This provides a description of how to use the argument.
func (arg ListOfStringsArgument) GetHelp() string {
	return arg.Help
}

// GetValue returns the current list of strings as an interface{}.
// It will return the actual value provided by the user or the default value if none was specified.
func (arg ListOfStringsArgument) GetValue() any {
	return *arg.Value
}

// SetValue sets the value of the ListOfStringsArgument.
// This is the list of strings provided by the user or set by default.
func (arg *ListOfStringsArgument) SetValue(value any) {
	*(arg.Value) = value.([]string)
}

// GetDefaultValue returns the default list of strings as an interface{}.
// This is used when the argument is not specified by the user.
func (arg ListOfStringsArgument) GetDefaultValue() any {
	return arg.DefaultValue
}

// ResetDefaultValue resets the value of the argument to the default value.
func (arg *ListOfStringsArgument) ResetDefaultValue() {
	*(arg.Value) = append(*(arg.Value), arg.DefaultValue...)
}

// IsRequired returns whether the argument is required.
// If true, the argument must be specified when running the program.
func (arg ListOfStringsArgument) IsRequired() bool {
	return arg.Required
}

// IsPresent checks if the argument was set in the command line.
func (arg ListOfStringsArgument) IsPresent() bool {
	return arg.Present
}

// Init initializes the ListOfStringsArgument with the provided values.
//
// The function sets the short name, long name, help message, value, and default value for the ListOfStringsArgument.
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
func (arg *ListOfStringsArgument) Init(value *[]string, shortName, longName string, defaultValue []string, required bool, help string) {
	arg.LongName, arg.ShortName = utils.GenerateLongAndShortNames(longName, shortName)

	arg.Required = required

	arg.Present = false

	arg.Help = help

	arg.Value = value

	arg.DefaultValue = defaultValue
}

// Consume processes the command-line arguments and sets the value of the ListOfStringsArgument.
//
// The function iterates through the provided arguments and checks if any of them match the short or long name of the ListOfStringsArgument.
// If a match is found, it sets the value of the ListOfStringsArgument to the next argument in the list and returns the remaining arguments.
// If the argument is required and not found, it returns an error.
//
// Parameters:
//   - arguments: A slice of strings representing the command-line arguments.
//
// Returns:
// - A slice of strings representing the remaining arguments after processing the ListOfStringsArgument.
func (arg *ListOfStringsArgument) Consume(arguments []string) ([]string, error) {
	// Initiate the value for the first time
	if arg.Value == nil {
		v := make([]string, 0)
		arg.Value = &v
	}

	sizeToConsume := 2

	if len(arguments) >= sizeToConsume {
		if (arguments[0] == arg.ShortName) || (arguments[0] == arg.LongName) {
			*arg.Value = append(*arg.Value, arguments[1])

			arg.Present = true

			return arguments[sizeToConsume:], nil
		}
	}

	return arguments, nil
}
