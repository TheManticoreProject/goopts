package arguments

// Argument defines the essential methods that all command-line argument types
// must implement. This interface allows different types of arguments (e.g., strings,
// integers, booleans) to be parsed and handled uniformly within the program.
type Argument interface {
	// GetShortName returns the short name (single-character) representation
	// of the argument (e.g., "-h" for help).
	GetShortName() string

	// GetLongName returns the long name (multi-character) representation
	// of the argument (e.g., "--help").
	GetLongName() string

	// GetHelp returns a help message or description for the argument.
	// This information is displayed when the user requests help or
	// usage information for the program.
	GetHelp() string

	// GetValue retrieves the current value of the argument after parsing.
	// This returns the actual value specified by the user during program execution.
	GetValue() any

	// GetDefaultValue returns the default value assigned to the argument.
	// This value is used if the user does not provide an explicit value for the argument.
	GetDefaultValue() any

	// IsRequired checks if the argument is marked as required.
	// If true, the program will enforce that this argument must be provided
	// by the user, otherwise an error will be thrown.
	IsRequired() bool

	// IsPresent checks if the argument was set in the command line.
	IsPresent() bool

	// Consume processes the command-line arguments, identifying and extracting
	// values that correspond to this specific argument. The method returns
	// a slice of the remaining unprocessed arguments after consuming the relevant ones.
	Consume(arguments []string) ([]string, error)
}
