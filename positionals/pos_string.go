package positionals

// StringPositionalArgument represents a positional command-line argument that expects a string value.
// This struct is used to define and handle a string input parameter, including its name, help message,
// and whether it is required.
//
// Fields:
//
//	Name (string): The name of the argument, used for display and reference purposes.
//	Help (string): A help message describing the purpose of the argument, shown in usage instructions.
//	Value (*string): A pointer to a string where the parsed value will be stored.
//	Required (bool): A flag indicating whether this argument must be provided. If set to true, the argument
//	                 is mandatory; otherwise, it is optional.
type StringPositionalArgument struct {
	Name     string
	Help     string  // Help message
	Value    *string // Values
	Required bool
}

// GetName retrieves the name of the positional argument.
//
// Returns:
//
//	(string): The name of the argument.
func (arg StringPositionalArgument) GetName() string {
	return arg.Name
}

// GetHelp retrieves the help message associated with the positional argument.
//
// Returns:
//
//	(string): The help message describing the argument.
func (arg StringPositionalArgument) GetHelp() string {
	return arg.Help
}

// GetValue retrieves the current value of the positional argument.
//
// Returns:
//
//	(any): The current value of the argument, or nil if not set.
func (arg StringPositionalArgument) GetValue() any {
	return *arg.Value
}

// IsRequired indicates whether the positional argument is mandatory.
//
// Returns:
//
//	(bool): True if the argument is required; otherwise, false.
func (arg StringPositionalArgument) IsRequired() bool {
	return arg.Required
}

// Init initializes the StringPositionalArgument with the provided values.
//
// The function sets the short name, long name, help message, value, and default value for the StringPositionalArgument.
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
func (arg *StringPositionalArgument) Init(value *string, name string, help string) {
	arg.Name = name

	arg.Help = help

	arg.Value = value

	arg.Required = true
}

// Consume processes the command-line arguments and sets the value of the StringPositionalArgument.
//
// The function iterates through the provided arguments and checks if any of them match the short or long name of the StringPositionalArgument.
// If a match is found, it sets the value of the StringPositionalArgument to the next argument in the list and returns the remaining arguments.
// If the argument is required and not found, it returns an error.
//
// Parameters:
//   - arguments: A slice of strings representing the command-line arguments.
//
// Returns:
// - A slice of strings representing the remaining arguments after processing the StringPositionalArgument.
func (arg StringPositionalArgument) Consume(arguments []string) ([]string, error) {
	sizeToConsume := 1

	if len(arguments) >= sizeToConsume {
		*arg.Value = arguments[0]

		return arguments[sizeToConsume:], nil
	}

	return arguments, nil
}
