package arguments

import (
	"fmt"

	"github.com/TheManticoreProject/goopts/utils"
)

// TcpPortArgument represents a command-line argument that specifies a TCP port.
// It holds information about the argument's short and long flag names, help message,
// and whether it is required or optional with a default value.
type TcpPortArgument struct {
	// ShortName is the short flag (e.g., "-p") used to specify the TCP port.
	// It can be empty if no short flag is defined.
	ShortName string
	// LongName is the long flag (e.g., "--port") used to specify the TCP port.
	// It can be empty if no long flag is defined.
	LongName string
	// Help provides a description of what this argument represents.
	// This message is displayed when showing help/usage information.
	Help string
	// Value stores the actual TCP port number provided by the user.
	// If no port is specified by the user, the Value will hold the DefaultValue.
	Value *int
	// DefaultValue is the port number to be used if the argument is not provided by the user.
	DefaultValue int
	// Required indicates whether this argument must be specified by the user.
	// If true, the argument must be included when running the program.
	Required bool
	// Present indicates whether this argument was set by the user during execution.
	// This can be used to differentiate between arguments that were provided and those that were not,
	// allowing for different handling of default values or other logic in the program.
	Present bool
}

// GetShortName returns the short flag name (e.g., "-p") of the TcpPortArgument.
// If the short name is not set, it returns an empty string.
func (arg TcpPortArgument) GetShortName() string {
	return arg.ShortName
}

// GetLongName returns the long flag name (e.g., "--port") of the TcpPortArgument.
// If the long name is not set, it returns an empty string.
func (arg TcpPortArgument) GetLongName() string {
	return arg.LongName
}

// GetHelp returns the help message associated with the TcpPortArgument.
// If the argument is optional, the default value is appended to the help message.
func (arg TcpPortArgument) GetHelp() string {
	if !arg.IsRequired() {
		return fmt.Sprintf("%s (default: %d)", arg.Help, arg.GetDefaultValue())
	} else {
		return arg.Help
	}
}

// GetValue retrieves the value of the TcpPortArgument.
// This is the port number provided by the user or set by default.
func (arg TcpPortArgument) GetValue() any {
	return *arg.Value
}

// SetValue sets the value of the TcpPortArgument.
// This is the port number provided by the user or set by default.
func (arg *TcpPortArgument) SetValue(value any) {
	*(arg.Value) = value.(int)
}

// GetDefaultValue returns the default value of the TcpPortArgument.
func (arg TcpPortArgument) GetDefaultValue() any {
	return arg.DefaultValue
}

// ResetDefaultValue resets the value of the argument to the default value.
func (arg *TcpPortArgument) ResetDefaultValue() {
	*(arg.Value) = arg.DefaultValue
}

// IsRequired checks if the TcpPortArgument is marked as required.
func (arg TcpPortArgument) IsRequired() bool {
	return arg.Required
}

// IsPresent checks if the argument was set in the command line.
func (arg TcpPortArgument) IsPresent() bool {
	return arg.Present
}

// Init initializes the TcpPortArgument with the provided values.
//
// Parameters:
//   - value (*int): Pointer to store the value of the argument.
//   - shortName (string): Short flag (e.g., "-p").
//   - longName (string): Long flag (e.g., "--port").
//   - defaultValue (int): Default value if the argument is not provided.
//   - required (bool): Indicates whether the argument is required.
//   - help (string): Help message to describe the argument.
func (arg *TcpPortArgument) Init(value *int, shortName, longName string, defaultValue int, required bool, help string) {
	arg.LongName, arg.ShortName = utils.GenerateLongAndShortNames(longName, shortName)

	arg.Required = required

	arg.Present = false

	arg.Help = help

	arg.Value = value

	arg.DefaultValue = defaultValue
}

// Consume processes the command-line arguments and sets the value of the TcpPortArgument.
//
// The function iterates through the provided arguments and checks if any of them match the short or long name of the TcpPortArgument.
// If a match is found, it sets the value of the TcpPortArgument to the next argument in the list and returns the remaining arguments.
// If the argument is required and not found, it returns an error.
//
// Parameters:
//   - arguments: A slice of strings representing the command-line arguments.
//
// Returns:
// - A slice of strings representing the remaining arguments after processing the TcpPortArgument.
func (arg *TcpPortArgument) Consume(arguments []string) ([]string, error) {
	sizeToConsume := 2

	if len(arguments) >= sizeToConsume {
		if (arguments[0] == arg.ShortName) || (arguments[0] == arg.LongName) {
			value, err := utils.StringToInt(arguments[1])
			if err != nil {
				// ParseInt failed, raise error and return the original arguments
				return arguments, fmt.Errorf("%s %s: could not parse integer: %s", arguments[0], arguments[1], err)
			}
			if value < 0 || value > 65535 {
				// ParseInt failed, raise error and return the original arguments
				return arguments, fmt.Errorf("%s %s: TCP port value has to be in range 0-65535", arguments[0], arguments[1])
			}
			(*arg.Value) = int(value)

			arg.Present = true

			return arguments[sizeToConsume:], nil
		}
	}

	return arguments, nil
}
