package positionals

// PositionalArgument is an interface that defines the common behavior for all positional argument types.
// Any struct that implements this interface can be used as a positional argument, allowing for a consistent
// approach to parsing and handling different argument types.
//
// Methods:
//
//	GetName() string: Retrieves the name of the argument, used for display and reference purposes.
//	GetHelp() string: Retrieves the help message describing the purpose of the argument.
//	GetValue() any: Retrieves the current value of the argument, which can be of any type.
//	IsRequired() bool: Indicates whether this argument is mandatory. Returns true if the argument must be provided.
//	Consume(arguments []string) ([]string, error): Parses and consumes the argument from a list of input strings. Returns the remaining arguments and an error if the parsing fails.
type PositionalArgument interface {
	// Retrieve the name
	GetName() string
	// Retrieve the help message
	GetHelp() string
	// Retrieve the value
	GetValue() any
	// Check if the argument is required
	IsRequired() bool
	// Parse the argument from the provided input
	Consume(arguments []string) ([]string, error)
}
