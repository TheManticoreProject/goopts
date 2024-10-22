package positionals

import "strings"

// BoolPositionalArgument represents a positional command-line argument that expects a boolean value.
// This struct is used to define and handle a boolean input parameter, including its name, help message,
// and whether it is required.
//
// Fields:
//
//	Name (string): The name of the argument, used for display and reference purposes.
//	Help (string): A help message describing the purpose of the argument, shown in usage instructions.
//	Value (*bool): A pointer to a boolean where the parsed value will be stored.
//	Required (bool): A flag indicating whether this argument must be provided. If set to true, the argument
//	                 is mandatory; otherwise, it is optional.
type BoolPositionalArgument struct {
	Name     string
	Help     string // Help message
	Value    *bool  // Values
	Required bool
}

// GetName retrieves the name of the positional argument.
//
// Returns:
//
//	(string): The name of the argument.
func (arg BoolPositionalArgument) GetName() string {
	return arg.Name
}

// GetHelp retrieves the help message associated with the positional argument.
//
// Returns:
//
//	(string): The help message describing the argument.
func (arg BoolPositionalArgument) GetHelp() string {
	return arg.Help
}

// GetValue retrieves the current value of the positional argument.
//
// Returns:
//
//	(any): The current boolean value of the argument, or nil if not set.
func (arg BoolPositionalArgument) GetValue() any {
	return *arg.Value
}

// IsRequired indicates whether the positional argument is mandatory.
//
// Returns:
//
//	(bool): Always returns true, as positional arguments are required.
func (arg BoolPositionalArgument) IsRequired() bool {
	// It is a positional argument, it is always required
	return true
}

// Init initializes the `BoolPositionalArgument` with a specified value, name, and help message.
//
// Parameters:
//
//	value (*bool): A pointer to the boolean that will hold the argument's value.
//	name (string): The name of the positional argument.
//	help (string): The help message describing the argument.
//
// Behavior:
//
//	This method sets up the `BoolPositionalArgument` and marks it as required.
func (arg *BoolPositionalArgument) Init(value *bool, name string, help string) {
	arg.Name = name

	arg.Help = help

	arg.Value = value

	arg.Required = true
}

// Consume processes the command-line arguments and sets the value of the BoolPositionalArgument.
//
// The function iterates through the provided arguments and checks if any of them match the short or long name of the BoolPositionalArgument.
// If a match is found, it sets the value of the BoolPositionalArgument to true and returns the remaining arguments.
// If the argument is required and not found, it returns an error.
//
// Parameters:
//   - arguments: A slice of strings representing the command-line arguments.
//
// Returns:
// - A slice of strings representing the remaining arguments after processing the BoolPositionalArgument.
func (arg BoolPositionalArgument) Consume(arguments []string) ([]string, error) {
	sizeToConsume := 1

	if len(arguments) >= sizeToConsume {
		if strings.EqualFold(strings.Trim(arguments[0], " "), "true") {
			*arg.Value = true
			return arguments[sizeToConsume:], nil
		} else if strings.EqualFold(strings.Trim(arguments[0], " "), "false") {
			*arg.Value = false
			return arguments[sizeToConsume:], nil
		}
	}

	return arguments, nil
}
