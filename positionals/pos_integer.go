package positionals

import "github.com/TheManticoreProject/goopts/utils"

// IntPositionalArgument represents a positional command-line argument that expects an integer value.
// This struct is used to define and handle an integer input parameter, including its name, help message,
// and whether it is required.
//
// Fields:
//
//	Name (string): The name of the argument, used for display and reference purposes.
//	Help (string): A help message describing the purpose of the argument, shown in usage instructions.
//	Value (*int): A pointer to an integer where the parsed value will be stored.
//	Required (bool): A flag indicating whether this argument must be provided. If set to true, the argument
//	                 is mandatory; otherwise, it is optional.
type IntPositionalArgument struct {
	Name     string
	Help     string // Help message
	Value    *int   // Values
	Required bool
}

// GetName retrieves the name of the positional argument.
//
// Returns:
//
//	(string): The name of the argument.
func (arg IntPositionalArgument) GetName() string {
	return arg.Name
}

// GetHelp retrieves the help message associated with the positional argument.
//
// Returns:
//
//	(string): The help message describing the argument.
func (arg IntPositionalArgument) GetHelp() string {
	return arg.Help
}

// GetValue retrieves the current value of the positional argument.
//
// Returns:
//
//	(any): The current integer value of the argument, or nil if not set.
func (arg IntPositionalArgument) GetValue() any {
	return *arg.Value
}

// IsRequired indicates whether the positional argument is mandatory.
//
// Returns:
//
//	(bool): True if the argument is required; otherwise, false.
func (arg IntPositionalArgument) IsRequired() bool {
	return arg.Required
}

// Init initializes the `IntPositionalArgument` with a specified value, name, and help message.
//
// Parameters:
//
//	value (*int): A pointer to the integer that will hold the argument's value.
//	name (string): The name of the positional argument.
//	help (string): The help message describing the argument.
//
// Behavior:
//
//	This method sets up the `IntPositionalArgument` and marks it as required.
func (arg *IntPositionalArgument) Init(value *int, name string, help string) {
	arg.Name = name

	arg.Help = help

	arg.Value = value

	arg.Required = true
}

// Consume processes the command-line arguments and sets the value of the IntPositionalArgument.
//
// The function iterates through the provided arguments and checks if any of them match the short or long name of the IntPositionalArgument.
// If a match is found, it sets the value of the IntPositionalArgument to the next argument in the list and returns the remaining arguments.
// If the argument is required and not found, it returns an error.
//
// Parameters:
//   - arguments: A slice of strings representing the command-line arguments.
//
// Returns:
// - A slice of strings representing the remaining arguments after processing the IntPositionalArgument.
func (arg IntPositionalArgument) Consume(arguments []string) ([]string, error) {
	sizeToConsume := 1

	if len(arguments) >= sizeToConsume {
		value, err := utils.StringToInt(arguments[0])
		if err != nil {
			// Return the original arguments if parsing fails
			return arguments, err
		}

		*arg.Value = int(value)

		return arguments[sizeToConsume:], nil
	}

	return arguments, nil
}
