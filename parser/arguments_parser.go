package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/p0dalirius/goopts/argumentgroup"
	"github.com/p0dalirius/goopts/arguments"
	"github.com/p0dalirius/goopts/positionals"
)

// ArgumentsParser is responsible for parsing command-line arguments for a program.
// It handles both positional and named arguments, supports argument groups, and provides
// usage information. The parser maintains a mapping of argument names to their corresponding
// values and manages required arguments.
//
// Fields:
//
//	Banner (string): A banner string that is printed before running the argument parsing.
//	                 This is typically used to display the program name or purpose.
//	PositionalArguments ([]*positionals.PositionalArgument): A slice of pointers to positional arguments
//	                                                         that the parser will manage.
//	DefaultGroup (*argumentgroup.ArgumentGroup): A pointer to the default argument group for organizing
//	                                              arguments, if applicable.
//	shortNameToArgument (map[string]arguments.Argument): A map that associates short flag names (e.g., "-v")
//	                                                      with their corresponding argument structures.
//	longNameToArgument (map[string]arguments.Argument): A map that associates long flag names (e.g., "--verbose")
//	                                                     with their corresponding argument structures.
//	requiredArguments ([]arguments.Argument): A slice of pointers to required arguments that must be provided
//	                                            by the user for the program to run successfully.
//	allArguments ([]arguments.Argument): A slice containing all arguments (both positional and named) that
//	                                       the parser manages.
//	Groups (map[string]*argumentgroup.ArgumentGroup): A map of named subgroups to organize related arguments
//	                                                  into logical categories for better structure and readability.
type ArgumentsParser struct {
	Banner           string
	ShowBannerOnHelp bool
	ShowBannerOnRun  bool

	PositionalArguments []positionals.PositionalArgument

	DefaultGroup *argumentgroup.ArgumentGroup

	Groups map[string]*argumentgroup.ArgumentGroup

	shortNameToArgument map[string]arguments.Argument
	longNameToArgument  map[string]arguments.Argument
	requiredArguments   []arguments.Argument
	allArguments        []arguments.Argument
}

// ArgumentIsPresent checks if a given argument is present in the parsed arguments.
// It supports both short (e.g., -e) and long (e.g., --example) argument names.
//
// Parameters:
//   - argumentName: The name of the argument to check. It should start with a dash ('-')
//     for short arguments or two dashes ('--') for long arguments.
//
// Returns:
// - bool: true if the argument is present; false otherwise.
//
// The function first verifies that the argument name has a valid length and starts
// with a dash. It then checks if the argument is in the `longNameToArgument` map
// for long arguments or `shortNameToArgument` map for short arguments, returning
// true if found and false if not.
func (ap *ArgumentsParser) ArgumentIsPresent(argumentName string) bool {
	if len(argumentName) < 2 || argumentName[0] != '-' {
		return false
	}

	if argumentName[1] == '-' {
		// Long argument flag (e.g., --example)
		if argument, exists := ap.longNameToArgument[argumentName]; exists {
			return argument.IsPresent()
		}
	} else {
		// Short argument flag (e.g., -e)
		if argument, exists := ap.shortNameToArgument[argumentName]; exists {
			return argument.IsPresent()
		}
	}

	// Argument not found
	return false
}

// populateMaps initializes the maps that store the associations between short and long argument names
// and their corresponding argument structures. This method is called to prepare the parser for argument
// handling and validation.
//
// Behavior:
//   - It creates or resets the `shortNameToArgument` and `longNameToArgument` maps to ensure they
//     are up-to-date.
//   - Registers arguments from the default argument group, storing their short and long names in the
//     respective maps. If an argument is required, it adds it to the `requiredArguments` slice.
//   - Registers arguments from all named subgroups in the `Groups` map, similarly storing their names
//     and tracking required arguments.
//
// Note:
//
//	This method should be called before parsing command-line arguments to ensure that all arguments
//	are correctly registered and accessible for lookup during the parsing process.
func (ap *ArgumentsParser) populateMaps() {
	// Create the maps even if they already exist
	ap.shortNameToArgument = make(map[string]arguments.Argument)
	ap.longNameToArgument = make(map[string]arguments.Argument)

	// Register arguments from the default group
	if ap.DefaultGroup != nil {
		for _, arg := range ap.DefaultGroup.Arguments {
			if shortName := arg.GetShortName(); shortName != "" {
				ap.shortNameToArgument[shortName] = arg
			}
			if longName := arg.GetLongName(); longName != "" {
				ap.longNameToArgument[longName] = arg
			}
			//
			if arg.IsRequired() {
				ap.requiredArguments = append(ap.requiredArguments, arg)
			}
			ap.allArguments = append(ap.allArguments, arg)
		}
	} else {
		// Not initialized, we init it for the first time
		ag := argumentgroup.ArgumentGroup{}
		ap.DefaultGroup = &ag
	}

	// Register arguments from all sub groups
	if ap.Groups != nil {
		for _, group := range ap.Groups {
			for _, arg := range group.Arguments {
				if shortName := arg.GetShortName(); shortName != "" {
					ap.shortNameToArgument[shortName] = arg
				}
				if longName := arg.GetLongName(); longName != "" {
					ap.longNameToArgument[longName] = arg
				}
				//
				if arg.IsRequired() {
					ap.requiredArguments = append(ap.requiredArguments, arg)
				}
				ap.allArguments = append(ap.allArguments, arg)
			}
		}
	} else {
		// Not initialized, we init it for the first time
		ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)
	}
}

// Parse processes the command-line arguments and sets the values for the defined arguments.
// This method handles both positional and named arguments, supports flags with values
// specified using "=", and checks for missing or unexpected arguments. It also provides
// a usage message if the "-h" or "--help" flags are present.
//
// Behavior:
//   - Populates maps for quick lookup of arguments based on their short and long names.
//   - Splits input arguments on "=" to allow for flags like "--key=value".
//   - Detects the presence of help flags ("-h" or "--help") and displays usage information.
//   - Separates positional arguments from named arguments based on the order of inputs.
//   - Validates that all required positional and named arguments are provided and parses them.
//   - Displays error messages for missing or extra arguments and exits if any errors are detected.
//
// Note:
//
//	This method terminates the program if it encounters errors or if help is requested.
//
// Example Usage:
//   - `./program positional1 positional2 --name=example`
//
// Errors:
//
//	If required arguments are missing or extra positional arguments are found, error messages
//	will be displayed, and the program will exit with a non-zero status.
func (ap *ArgumentsParser) ParseFrom(index int) {
	debug := false

	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Calling populateMaps()\n", index)
	}
	ap.populateMaps()

	if len(ap.Banner) != 0 && ap.ShowBannerOnRun {
		fmt.Printf("%s\n\n", ap.Banner)
	}

	// Store error messages of parsing arguments
	errorMessages := []string{}

	// Prepare arguments and split on "=" for `--arg=value`
	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Prepare arguments and split on \"=\" for `--arg=value`\n", index)
	}
	arguments := []string{}
	for _, arg := range os.Args[index:] {
		if strings.Contains(arg, "=") && strings.HasPrefix(arg, "-") {
			arguments = append(arguments, strings.SplitN(arg, "=", 2)...)
		} else {
			arguments = append(arguments, arg)
		}
	}

	// Check if -h or --help are present
	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Check if -h or --help are present\n", index)
	}
	if slices.Contains(arguments, "-h") || slices.Contains(arguments, "--help") {
		ap.Usage()
		os.Exit(0)
	}

	// Reset all arguments to their default values
	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Reset all arguments to their default values\n", index)
	}
	for _, arg := range ap.allArguments {
		arg.ResetDefaultValue()
	}

	// Split between positional arguments and other arguments
	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Split between positional arguments and other arguments\n", index)
	}
	potentialPositionalArguments := []string{}
	otherArguments := []string{}
	parsingPositionalArguments := true
	for _, arg := range arguments {
		if strings.HasPrefix(arg, "-") {
			parsingPositionalArguments = false
		}
		if parsingPositionalArguments {
			potentialPositionalArguments = append(potentialPositionalArguments, arg)
		} else {
			otherArguments = append(otherArguments, arg)
		}
	}

	// Parse the positional arguments first
	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Parse the positional arguments first\n", index)
	}
	missingPositionalArguments := []string{}
	for k, posarg := range ap.PositionalArguments {
		if k < len(potentialPositionalArguments) {
			posarg.Consume([]string{potentialPositionalArguments[k]})
		} else {
			missingPositionalArguments = append(missingPositionalArguments, posarg.GetName())
		}
	}
	if len(missingPositionalArguments) != 0 {
		if len(missingPositionalArguments) == 1 {
			errorMessages = append(errorMessages, fmt.Sprintf("Missing %d positional argument: <%s>.", len(missingPositionalArguments), missingPositionalArguments[0]))
		} else {
			errmsg := fmt.Sprintf("Missing %d positional arguments:", len(missingPositionalArguments))
			for _, posarg := range missingPositionalArguments {
				errmsg = errmsg + fmt.Sprintf(" <%s>", posarg)
			}
			errmsg = errmsg + "."
			errorMessages = append(errorMessages, errmsg)
		}
	}
	if len(potentialPositionalArguments) > len(ap.PositionalArguments) {
		leftoverPositionalArguments := potentialPositionalArguments[len(ap.PositionalArguments):]
		if len(leftoverPositionalArguments) == 1 {
			errorMessages = append(errorMessages, fmt.Sprintf("Got %d more positional argument than expected: \"%s\".", len(leftoverPositionalArguments), leftoverPositionalArguments[0]))
		} else {
			errmsg := fmt.Sprintf("Got %d more positional argument than expected: ", len(leftoverPositionalArguments))
			for _, loposarg := range leftoverPositionalArguments {
				errmsg = errmsg + fmt.Sprintf(" \"%s\"", loposarg)
			}
			errmsg = errmsg + "."
			errorMessages = append(errorMessages, errmsg)
		}
	}

	// Parse all other arguments
	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Parse all other arguments\n", index)
	}
	for k, otherarg := range otherArguments {
		if strings.HasPrefix(otherarg, "--") {
			// Long flag name
			if debug {
				fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] | Processing argument with long name: [%s]\n", index, otherarg)
			}
			if _, exists := ap.longNameToArgument[otherarg]; exists {
				if debug {
					fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] |  | Argument exists, consuming argument.\n", index)
				}
				arg := ap.longNameToArgument[otherarg]
				_, err := arg.Consume(otherArguments[k:])
				if err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Error parsing argument: %s", err))
				}
			}
		} else if strings.HasPrefix(otherarg, "-") {
			// Short flag name
			if debug {
				fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] | Processing argument with short name: [%s]\n", index, otherarg)
			}
			if _, exists := ap.shortNameToArgument[otherarg]; exists {
				if debug {
					fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] |  | Argument exists, consuming argument.\n", index)
				}
				arg := ap.shortNameToArgument[otherarg]
				_, err := arg.Consume(otherArguments[k:])
				if err != nil {
					errorMessages = append(errorMessages, fmt.Sprintf("Error parsing argument: %s", err))
				}
			}
		}
	}

	// Check if all required arguments have been parsed
	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Check if all required arguments have been parsed\n", index)
	}
	requiredArgumentsMissing := []string{}
	for _, arg := range ap.requiredArguments {
		if !arg.IsPresent() {
			requiredArgumentsMissing = append(requiredArgumentsMissing, arg.GetLongName())
		}
	}
	if len(requiredArgumentsMissing) != 0 {
		if len(requiredArgumentsMissing) == 1 {
			errorMessages = append(errorMessages, fmt.Sprintf("Missing required argument \"%s\"", requiredArgumentsMissing[0]))
		} else {
			errorMessages = append(errorMessages, fmt.Sprintf("Missing required arguments \"%s\"", strings.Join(requiredArgumentsMissing, "\", \"")))
		}
	}

	// Check if all required arguments in groups have been parsed
	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Check if all required arguments in groups have been parsed\n", index)
	}
	for _, group := range ap.Groups {
		argumentsPresent := []string{}
		argumentsMissing := []string{}
		for _, arg := range group.Arguments {
			if arg.IsPresent() {
				argumentsPresent = append(argumentsPresent, arg.GetLongName())
			} else {
				argumentsMissing = append(argumentsMissing, arg.GetLongName())
			}
		}

		if group.Type == argumentgroup.ARGUMENT_GROUP_TYPE_REQUIRED_MUTUALLY_EXCLUSIVE {
			// One needs to be set, and one only
			if len(argumentsPresent) == 0 {
				if len(argumentsMissing) == 1 {
					errorMessages = append(errorMessages, fmt.Sprintf("The argument \"%s\" needs to be set.", argumentsMissing[0]))
				} else if len(argumentsMissing) > 1 {
					errorMessages = append(errorMessages, fmt.Sprintf("At least one of the arguments \"%s\" needs to be set.", strings.Join(argumentsMissing, "\", \"")))
				}
			} else if len(argumentsPresent) > 1 {
				errorMessages = append(errorMessages, fmt.Sprintf("Arguments \"%s\" cannot be set together.", strings.Join(argumentsPresent, "\", \"")))
			}
		} else if group.Type == argumentgroup.ARGUMENT_GROUP_TYPE_NOT_REQUIRED_MUTUALLY_EXCLUSIVE {
			// None can be set but if one is set then only one has to be set
			if len(argumentsPresent) > 1 {
				errorMessages = append(errorMessages, fmt.Sprintf("Arguments \"%s\" cannot be set together.", strings.Join(argumentsPresent, "\", \"")))
			}
		} else if group.Type == argumentgroup.ARGUMENT_GROUP_TYPE_DEPENDENT {
			// If one is set, all need to be set
			if len(argumentsMissing) != 0 {
				if len(argumentsPresent) > 1 {
					errorMessages = append(errorMessages, fmt.Sprintf("When arguments \"%s\" are set, \"%s\" need to be set too.", strings.Join(argumentsPresent, "\", \""), strings.Join(argumentsMissing, "\", \"")))
				} else if len(argumentsPresent) == 1 {
					errorMessages = append(errorMessages, fmt.Sprintf("When argument \"%s\" is set, \"%s\" need to be set too.", argumentsPresent[0], strings.Join(argumentsMissing, "\", \"")))
				}
			}
		}
	}

	if len(errorMessages) != 0 {
		if debug {
			fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Error messages found, printing usage and exiting\n", index)
		}
		ap.UsageFrom(index)
		for _, errmsg := range errorMessages {
			fmt.Printf("[!] %s\n", errmsg)
		}
		os.Exit(1)
	}

	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] No errors found, arguments are parsed!\n", index)
	}
}

// Parse parses the arguments and returns the parsed arguments.
//
// Returns:
// - A map of parsed arguments.
func (ap *ArgumentsParser) Parse() {
	// We start parsing from index 1 because the first argument (0) is the program name
	ap.ParseFrom(1)
}

// Get retrieves the value of the argument specified by its short or long name.
//
// Parameters:
//   - name: The name of the argument to retrieve. This can be either the short name (single character)
//     or the long name (string) of the argument.
//
// Returns:
// - The value of the argument as an interface{} if the argument is found.
// - An error if the argument maps are not initialized or if the argument with the specified name is not found.
//
// The function first checks if the argument maps (shortNameToArgument and longNameToArgument) are initialized.
// If they are not, it returns an error indicating that the argument maps are not initialized.
// It then looks up the argument by its short name in the shortNameToArgument map.
// If the argument is not found, it looks up the argument by its long name in the longNameToArgument map.
// If the argument is still not found, it returns an error indicating that the argument with the specified name is not found.
// If the argument is found, it returns the value of the argument.
func (ap *ArgumentsParser) Get(argumentFlag string) (interface{}, error) {
	if ap.shortNameToArgument == nil || ap.longNameToArgument == nil {
		return nil, fmt.Errorf("argument maps are not initialized")
	}

	var arg arguments.Argument
	var exists bool

	if arg, exists = ap.shortNameToArgument[argumentFlag]; !exists {
		if arg, exists = ap.longNameToArgument[argumentFlag]; !exists {
			return nil, fmt.Errorf("argument '%s' not found", argumentFlag)
		}
	}

	return arg.GetValue(), nil
}

// Usage prints the usage information for the command-line arguments.
//
// The function first prints the banner, followed by the usage string, which includes the name of the executable.
// It then iterates through all the arguments in the DefaultGroup and prints their short name, long name, and help description.
// The arguments are formatted using the padding format calculated by the computePaddingFormat function.
//
// After printing the arguments in the DefaultGroup, the function iterates over the named subgroups in the Groups map.
// For each subgroup, it prints the group name and the arguments within that group, including their short name, long name, and help description.
// The subgroup arguments are also formatted using the padding format calculated by the computePaddingFormat function.
//
// The function ensures that the usage information is displayed in a clear and organized manner, making it easy for users to understand
// the available command-line arguments and their descriptions.
func (ap *ArgumentsParser) UsageFrom(index int) {
	// Create usage string
	usage := filepath.Base(os.Args[0])
	if index != 0 {
		for k := range index - 1 {
			usage += " " + os.Args[k+1]
		}
	}

	// Append positional arguments
	for _, posarg := range ap.PositionalArguments {
		usage += fmt.Sprintf(" <%s>", posarg.GetName())
	}

	// Append default group arguments
	for _, argument := range ap.DefaultGroup.Arguments {
		output := generateArgumentForUsageLine(argument)
		if len(output) == 0 {
			// This should not happen
		} else {
			usage += " " + output
		}
	}

	// Prepare to iterate over groups sorted alphabetically by their names
	groupNames := make([]string, 0, len(ap.Groups))
	for groupname := range ap.Groups {
		groupNames = append(groupNames, groupname)
	}
	sort.Strings(groupNames)
	// Append arguments in groups
	for _, groupname := range groupNames {
		group := ap.Groups[groupname]
		for _, argument := range group.Arguments {
			output := generateArgumentForUsageLine(argument)
			if len(output) == 0 {
				// This should not happen
			} else {
				usage += " " + output
			}
		}
	}

	// Starting to print usage
	fmt.Printf("Usage: %s\n", usage)
	if len(ap.DefaultGroup.Arguments) != 0 {
		fmt.Printf("\n")
		fmtString := fmt.Sprintf("  %s %%s\n", computePaddingFormat(ap.DefaultGroup.Arguments))
		for _, argument := range ap.DefaultGroup.Arguments {
			printArgumentLineInHelp(argument, fmtString)
		}
	}

	// Iterate over groups
	for _, groupname := range groupNames {
		group := ap.Groups[groupname]

		fmt.Printf("\n  %s:\n", group.Name)
		fmtString := fmt.Sprintf("    %s %%s\n", computePaddingFormat(group.Arguments))
		for _, argument := range group.Arguments {
			printArgumentLineInHelp(argument, fmtString)
		}
	}

	fmt.Printf("\n")
}

// Usage prints the usage information for the command-line arguments.
func (ap *ArgumentsParser) Usage() {
	// We start printing from index 0 because the first argument (0)
	// is the program name and we need it in the usage
	ap.UsageFrom(0)
}

// generateArgumentForUsageLine generates a formatted string representing a command-line argument
// for inclusion in the usage line of a help message.
//
// Parameters:
//
//	arg (arguments.Argument): The argument to be formatted.
//
// Returns:
//
//	(string): A string describing the argument, including its type (if applicable) and name,
//	          enclosed in square brackets if the argument is optional.
//
// Behavior:
//   - Determines the argument's type (e.g., `StringArgument`, `IntArgument`, `TcpPortArgument`, etc.)
//     and appends an appropriate type descriptor (e.g., "<string>", "<int>", "<tcp port>").
//   - Checks if the argument has a long or short name and uses it to build the output.
//   - If the argument is not required, encloses the output string in square brackets.
//   - Logs an error message if the argument type is not recognized.
func generateArgumentForUsageLine(arg arguments.Argument) string {
	output := ""

	shortName := arg.GetShortName()
	longName := arg.GetLongName()

	// Don't know why I need to remove the * only for this one, it should be argument.(*arguments.BoolArgument)
	if _, ok := arg.(*arguments.BoolArgument); ok {
		if len(longName) != 0 {
			output = longName
		} else if len(shortName) != 0 {
			output = shortName
		}
	} else if _, ok := arg.(*arguments.StringArgument); ok {
		if len(longName) != 0 {
			output = fmt.Sprintf("%s <string>", longName)
		} else if len(shortName) != 0 {
			output = fmt.Sprintf("%s <string>", shortName)
		}
	} else if _, ok := arg.(*arguments.IntArgument); ok {
		if len(longName) != 0 {
			output = fmt.Sprintf("%s <int>", longName)
		} else if len(shortName) != 0 {
			output = fmt.Sprintf("%s <int>", shortName)
		}
	} else if _, ok := arg.(*arguments.IntRangeArgument); ok {
		if len(longName) != 0 {
			output = fmt.Sprintf("%s <int>", longName)
		} else if len(shortName) != 0 {
			output = fmt.Sprintf("%s <int>", shortName)
		}
	} else if _, ok := arg.(*arguments.TcpPortArgument); ok {
		if len(longName) != 0 {
			output = fmt.Sprintf("%s <tcp port>", longName)
		} else if len(shortName) != 0 {
			output = fmt.Sprintf("%s <tcp port>", shortName)
		}
	} else if _, ok := arg.(*arguments.MapOfHttpHeadersArgument); ok {
		if len(longName) != 0 {
			output = fmt.Sprintf("%s <http header>", longName)
		} else if len(shortName) != 0 {
			output = fmt.Sprintf("%s <http header>", shortName)
		}
	} else if _, ok := arg.(*arguments.ListOfIntsArgument); ok {
		if len(longName) != 0 {
			output = fmt.Sprintf("%s <int>", longName)
		} else if len(shortName) != 0 {
			output = fmt.Sprintf("%s <int>", shortName)
		}
	} else if _, ok := arg.(*arguments.ListOfStringsArgument); ok {
		if len(longName) != 0 {
			output = fmt.Sprintf("%s <string>", longName)
		} else if len(shortName) != 0 {
			output = fmt.Sprintf("%s <string>", shortName)
		}
	}

	if !arg.IsRequired() {
		output = fmt.Sprintf("[%s]", output)
	}

	return output
}

// printArgumentLineInHelp formats and prints a line in the help message for a given command-line argument.
//
// Parameters:
//
//	arg (arguments.Argument): The argument to be described in the help output.
//	fmtString (string): A format string specifying how to print the argument line.
//	                    It typically includes placeholders for the argument flags and the help message.
//
// Behavior:
//   - Retrieves the short and long names of the argument, then combines them into a flags string.
//   - Depending on the argument type, appends additional information to the flags string to indicate the expected type
//     (e.g., "<string>" or "<int>").
//   - If the argument is of type `BoolArgument`, the help message includes its default value.
//   - Outputs the formatted argument line using the provided format string.
func printArgumentLineInHelp(arg arguments.Argument, fmtString string) {
	shortName := arg.GetShortName()
	longName := arg.GetLongName()
	help := arg.GetHelp()

	// Prepare flags format
	flags_string := ""
	if (len(shortName) != 0) && (len(longName) != 0) {
		flags_string = fmt.Sprintf("%s, %s", shortName, longName)
	} else if (len(shortName) == 0) && (len(longName) != 0) {
		flags_string = longName
	} else if (len(shortName) != 0) && (len(longName) == 0) {
		flags_string = shortName
	}

	// Update flags with types
	if _, ok := arg.(*arguments.BoolArgument); ok {
		help = fmt.Sprintf("%s (default: %v)", help, arg.GetDefaultValue())
	} else if _, ok := arg.(*arguments.StringArgument); ok {
		flags_string = flags_string + " <string>"
	} else if _, ok := arg.(*arguments.ListOfStringsArgument); ok {
		flags_string = flags_string + " <string>"
	} else if _, ok := arg.(*arguments.IntArgument); ok {
		flags_string = flags_string + " <int>"
	} else if _, ok := arg.(*arguments.IntRangeArgument); ok {
		flags_string = flags_string + " <int>"
	} else if _, ok := arg.(*arguments.ListOfIntsArgument); ok {
		flags_string = flags_string + " <int>"
	} else if _, ok := arg.(*arguments.TcpPortArgument); ok {
		flags_string = flags_string + " <tcp port>"
	} else if _, ok := arg.(*arguments.MapOfHttpHeadersArgument); ok {
		flags_string = flags_string + " <http header>"
	}

	// Print the line
	fmt.Printf(fmtString, flags_string, help)
}

// PrintArgumentTree prints the argument tree for the ArgumentsParser.
//
// The function prints the banner, the list of arguments, and the subgroups in a tree-like structure.
// Each level of indentation is represented by "  │ ". The output includes the banner, the names of the arguments,
// and the names of the subgroups within the ArgumentsParser.
func (ap *ArgumentsParser) PrintArgumentTree() {
	indent := 0
	//
	fmt.Printf("<ArgumentsParser>\n")

	fmt.Printf("  ├─ Banner: \"%s\"\n", ap.Banner)

	fmt.Printf("  ├─ Arguments (%d): \n", len(ap.DefaultGroup.Arguments))
	for _, argument := range ap.DefaultGroup.Arguments {
		argtype := ""
		if _, ok := argument.(*arguments.BoolArgument); ok {
			argtype = "bool"
		} else if _, ok := argument.(*arguments.StringArgument); ok {
			argtype = "string"
		} else if _, ok := argument.(*arguments.IntArgument); ok {
			argtype = "int"
		}
		fmt.Printf("  │   │   ├─ (\"%s\",\"%s\") [%s] \"%s\"\n", argument.GetShortName(), argument.GetLongName(), argtype, argument.GetHelp())
	}
	fmt.Printf("  │   │   └──\n")
	fmt.Printf("  │   └──\n")

	// Call subgroups
	fmt.Printf("  ├─ SubGroups: \n")
	for _, subgroup := range ap.Groups {
		subgroup.PrintArgumentTree(indent + 1)
	}
	fmt.Printf("  └──\n")
}

// computePaddingFormat calculates the format string for padding argument flags.
//
// The function iterates through all the arguments in the DefaultGroup and determines the maximum length
// of the combined short and long name strings. It then creates a format string that can be used to
// align the argument flags in the usage output.
//
// Returns:
// - A format string that can be used to pad the argument flags to the maximum length.
func computePaddingFormat(args []arguments.Argument) string {
	max_len_flags_string := 15
	for _, argument := range args {
		shortName := argument.GetShortName()
		longName := argument.GetLongName()

		flags_string := ""
		if (len(shortName) != 0) && (len(longName) != 0) {
			flags_string = fmt.Sprintf("%s, %s", shortName, longName)
		} else if (len(shortName) == 0) && (len(longName) != 0) {
			flags_string = longName
		} else if (len(shortName) != 0) && (len(longName) == 0) {
			flags_string = shortName
		}

		arg := argument
		if _, ok := arg.(*arguments.StringArgument); ok {
			flags_string = flags_string + " <string>"
		} else if _, ok := arg.(*arguments.ListOfStringsArgument); ok {
			flags_string = flags_string + " <string>"
		} else if _, ok := arg.(*arguments.IntArgument); ok {
			flags_string = flags_string + " <int>"
		} else if _, ok := arg.(*arguments.IntRangeArgument); ok {
			flags_string = flags_string + " <int>"
		} else if _, ok := arg.(*arguments.ListOfIntsArgument); ok {
			flags_string = flags_string + " <int>"
		} else if _, ok := arg.(*arguments.TcpPortArgument); ok {
			flags_string = flags_string + " <tcp port>"
		} else if _, ok := arg.(*arguments.MapOfHttpHeadersArgument); ok {
			flags_string = flags_string + " <http header>"
		}

		if len(flags_string) > max_len_flags_string {
			max_len_flags_string = len(flags_string)
		}
	}

	return fmt.Sprintf("%%-%ds", max_len_flags_string)
}
