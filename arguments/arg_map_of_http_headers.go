package arguments

import (
	"strings"

	"github.com/p0dalirius/goopts/utils"
)

// MapOfHttpHeadersArgument represents a command-line argument that expects a map of HTTP headers.
// It contains information about the argument's short and long flag names, help message,
// the default value, and whether the argument is required.
type MapOfHttpHeadersArgument struct {
	// ShortName is the short flag (e.g., "-H") used to specify the HTTP headers.
	// It can be empty if no short flag is defined.
	ShortName string
	// LongName is the long flag (e.g., "--headers") used to specify the HTTP headers.
	// It can be empty if no long flag is defined.
	LongName string
	// Help provides a description of what this argument represents.
	// This message is displayed when showing help/usage information.
	Help string
	// Value stores the actual map of HTTP headers provided by the user.
	// If no value is specified by the user, Value will hold the DefaultValue.
	Value *map[string]string
	// DefaultValue is the map of HTTP headers to be used if the argument is not provided by the user.
	DefaultValue map[string]string
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
func (arg MapOfHttpHeadersArgument) GetShortName() string {
	return arg.ShortName
}

// GetLongName returns the long flag name of the argument.
// If no long flag is defined, it returns an empty string.
func (arg MapOfHttpHeadersArgument) GetLongName() string {
	return arg.LongName
}

// GetHelp returns the help message of the argument.
// This provides a description of how to use the argument.
func (arg MapOfHttpHeadersArgument) GetHelp() string {
	return arg.Help
}

// GetValue returns the current map of HTTP headers as an interface{}.
// It will return the actual value provided by the user or the default value if none was specified.
func (arg MapOfHttpHeadersArgument) GetValue() any {
	return *arg.Value
}

// GetDefaultValue returns the default map of HTTP headers as an interface{}.
// This is used when the argument is not specified by the user.
func (arg MapOfHttpHeadersArgument) GetDefaultValue() any {
	return arg.DefaultValue
}

// IsRequired returns whether the argument is required.
// If true, the argument must be specified when running the program.
func (arg MapOfHttpHeadersArgument) IsRequired() bool {
	return arg.Required
}

// IsPresent checks if the argument was set in the command line.
func (arg MapOfHttpHeadersArgument) IsPresent() bool {
	return arg.Present
}

// Init initializes the MapOfHttpHeadersArgument with the provided values.
//
// The function sets the short name, long name, help message, value, and default value for the MapOfHttpHeadersArgument.
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
func (arg *MapOfHttpHeadersArgument) Init(value *map[string]string, shortName, longName string, defaultValue map[string]string, required bool, help string) {
	if len(shortName) == 0 {
		arg.ShortName = ""
	} else {
		arg.ShortName = "-" + utils.StripLeftDashes(shortName)
	}

	if len(longName) == 0 {
		arg.LongName = ""
	} else {
		arg.LongName = "--" + utils.StripLeftDashes(longName)
	}

	arg.Required = required

	arg.Present = false

	arg.Help = help

	arg.Value = value
	if (*arg.Value) == nil {
		(*arg.Value) = make(map[string]string, 0)
	}
	for key, value := range defaultValue {
		(*arg.Value)[key] = value
	}

	arg.DefaultValue = defaultValue
}

// Consume processes the command-line arguments and sets the value of the MapOfHttpHeadersArgument.
//
// The function iterates through the provided arguments and checks if any of them match the short or long name of the MapOfHttpHeadersArgument.
// If a match is found, it sets the value of the MapOfHttpHeadersArgument to the next argument in the list and returns the remaining arguments.
// If the argument is required and not found, it returns an error.
//
// Parameters:
//   - arguments: A slice of strings representing the command-line arguments.
//
// Returns:
// - A slice of strings representing the remaining arguments after processing the MapOfHttpHeadersArgument.
func (arg *MapOfHttpHeadersArgument) Consume(arguments []string) ([]string, error) {
	// Initiate the value for the first time
	if (*arg.Value) == nil {
		(*arg.Value) = make(map[string]string, 0)
	}

	sizeToConsume := 2

	if len(arguments) >= sizeToConsume {
		if (arguments[0] == arg.ShortName) || (arguments[0] == arg.LongName) {
			if strings.Contains(arguments[1], ":") {
				http_headers := strings.SplitN(arguments[1], ":", 2)
				(*arg.Value)[strings.Trim(http_headers[0], " ")] = strings.Trim(http_headers[1], " ")
			} else {
				(*arg.Value)[strings.Trim(arguments[1], " ")] = ""
			}

			arg.Present = true

			return arguments[sizeToConsume:], nil
		}
	}

	return arguments, nil
}
