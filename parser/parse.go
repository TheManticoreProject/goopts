package parser

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/TheManticoreProject/goopts/argumentgroup"
	"github.com/TheManticoreProject/goopts/arguments"
	"github.com/TheManticoreProject/goopts/positionals"
)

// isMutuallyExclusiveGroup reports whether a group type enforces mutual exclusivity between its
// members, in which case the group rule alone decides whether a member has to be set.
//
// Parameters:
//   - groupType: The type of the argument group.
//
// Returns:
// - true for the required and not required mutually exclusive group types, false otherwise.
func isMutuallyExclusiveGroup(groupType int) bool {
	return groupType == argumentgroup.ARGUMENT_GROUP_TYPE_REQUIRED_MUTUALLY_EXCLUSIVE ||
		groupType == argumentgroup.ARGUMENT_GROUP_TYPE_NOT_REQUIRED_MUTUALLY_EXCLUSIVE
}

// sortedGroupNames returns the names of all argument groups sorted alphabetically, including the
// empty name of the default group.
//
// Iterating groups through this method instead of ranging over the `Groups` map keeps every output
// derived from that iteration reproducible, since Go randomizes map iteration order. It also matches
// the alphabetical group ordering used by the usage message.
//
// Returns:
// - A slice of group names sorted alphabetically.
func (ap *ArgumentsParser) sortedGroupNames() []string {
	groupNames := make([]string, 0, len(ap.Groups))
	for groupName := range ap.Groups {
		groupNames = append(groupNames, groupName)
	}
	slices.Sort(groupNames)

	return groupNames
}

// formatGroupErrorMessage builds the error message reported for an unsatisfied argument group
// constraint, prefixed with the name of the group the constraint belongs to.
//
// A parser can declare several constraint groups, each producing a message of the same shape, so
// naming the group is what makes those lines distinguishable and lets them be matched with the
// named sections of the usage message.
//
// Parameters:
//   - groupName: The name of the argument group the constraint belongs to. It is empty for the
//     default group.
//   - message: The constraint message, written in lower case so it reads correctly after the
//     group name.
//
// Returns:
//   - The message prefixed with the group name, or the message alone with its first letter
//     capitalized when the group has no name.
func formatGroupErrorMessage(groupName, message string) string {
	if len(groupName) == 0 {
		if len(message) == 0 {
			return message
		}

		return strings.ToUpper(message[:1]) + message[1:]
	}

	return fmt.Sprintf("%s: %s", groupName, message)
}

// populateMaps initializes the maps that store the associations between short and long argument names
// and their corresponding argument structures. This method is called to prepare the parser for argument
// handling and validation.
//
// Behavior:
//   - It creates or resets the `shortNameToArgument` and `longNameToArgument` maps to ensure they
//     are up-to-date.
//   - Registers arguments from the default argument group, storing their short and long names in the
//     respective maps. If an argument is required, it adds it to the `requiredArguments` slice, unless
//     it belongs to a mutually exclusive group, where the group rule decides whether it has to be set.
//   - Registers arguments from all named subgroups in the `Groups` map, similarly storing their names
//     and tracking required arguments.
//   - Initializes the `ParsedArguments` maps of `parsingState`, which is the state parsing results are
//     recorded into and is not necessarily the parser's own `ParsingState`.
//
// Parameters:
//   - parsingState: The parsing state that will be populated during the parse.
//
// Note:
//
//	This method should be called before parsing command-line arguments to ensure that all arguments
//	are correctly registered and accessible for lookup during the parsing process.
func (ap *ArgumentsParser) populateMaps(parsingState *ParsingState) {
	// Reset lookup maps and per-parse slices so repeated calls stay consistent
	ap.shortNameToArgument = make(map[string]arguments.Argument)
	ap.longNameToArgument = make(map[string]arguments.Argument)
	ap.requiredArguments = ap.requiredArguments[:0]
	ap.allArguments = ap.allArguments[:0]

	if ap.Groups == nil {
		ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)
	}

	// Ensure the default group exists so downstream code can rely on it
	if ap.Groups[""] == nil {
		ap.Groups[""] = &argumentgroup.ArgumentGroup{}
	}

	// Register arguments from every group exactly once (the default group is keyed by ""),
	// in a stable order so that error messages built from these slices are reproducible
	for _, groupName := range ap.sortedGroupNames() {
		group := ap.Groups[groupName]
		for _, arg := range group.Arguments {
			if shortName := arg.GetShortName(); shortName != "" {
				ap.shortNameToArgument[shortName] = arg
			}
			if longName := arg.GetLongName(); longName != "" {
				ap.longNameToArgument[longName] = arg
			}
			// In a mutually exclusive group, whether a member has to be set is decided by the
			// group rule, so its individual required flag is not enforced on its own: doing so
			// would contradict the group and make every other member unusable
			if arg.IsRequired() && !isMutuallyExclusiveGroup(group.Type) {
				ap.requiredArguments = append(ap.requiredArguments, arg)
			}
			ap.allArguments = append(ap.allArguments, arg)
		}
	}

	// Initialize the maps of the parsing state that will be written to during this parse,
	// which is not necessarily the parser's own ap.ParsingState
	if parsingState.ParsedArguments.PositionalArguments == nil {
		parsingState.ParsedArguments.PositionalArguments = make(map[string]*positionals.PositionalArgument)
	}
	if parsingState.ParsedArguments.LongNameToArgument == nil {
		parsingState.ParsedArguments.LongNameToArgument = make(map[string]*arguments.Argument)
	}
	if parsingState.ParsedArguments.ShortNameToArgument == nil {
		parsingState.ParsedArguments.ShortNameToArgument = make(map[string]*arguments.Argument)
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
//   - Reports an error for any argument starting with "-" that matches no registered short or long name.
//   - Displays error messages for missing, unknown or extra arguments and exits if any errors are detected.
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
//	If required arguments are missing, unknown flags are supplied, or extra positional arguments
//	are found, error messages will be displayed, and the program will exit with a non-zero status.
func (ap *ArgumentsParser) ParseFrom(index int, parsingState *ParsingState) {
	ap.populateMaps(parsingState)

	// Print the banner if it is set and the option is enabled
	if len(ap.Banner) != 0 && ap.Options.ShowBannerOnRun {
		fmt.Printf("%s\n\n", ap.Banner)
	}

	// Handle subparsers if enabled
	if ap.SubParsers.Enabled && len(ap.SubParsers.Parsers) != 0 {
		if index < len(parsingState.RawArguments) {
			subparserName := parsingState.RawArguments[index]
			if subparserName == "-h" || subparserName == "--help" {
				ap.UsageFrom(index, parsingState)
				os.Exit(0)
			}
			lookupName := subparserName
			if ap.SubParsers.CaseInsensitive {
				lookupName = strings.ToLower(subparserName)
			}
			if asp, exists := ap.SubParsers.Parsers[lookupName]; exists {
				// Set the subparser name value to the pointer
				*(ap.SubParsers.Value) = lookupName
				asp.ParseFrom(index+1, parsingState)
				return
			} else {
				parsingState.AddErrorMessage(fmt.Sprintf("No subparser with name \"%s\" was found.", lookupName))
			}
		} else {
			ap.UsageFrom(index, parsingState)
			os.Exit(1)
		}
	} else {
		// Prepare arguments and split on "=" for `--arg=value`
		// The index comes from the caller and can point past the end of the raw arguments, in
		// which case there is simply nothing left to parse
		rawArguments := []string{}
		if index >= 0 && index < len(parsingState.RawArguments) {
			rawArguments = parsingState.RawArguments[index:]
		}

		arguments := []string{}
		for _, arg := range rawArguments {
			if strings.Contains(arg, "=") && strings.HasPrefix(arg, "-") {
				arguments = append(arguments, strings.SplitN(arg, "=", 2)...)
			} else {
				arguments = append(arguments, arg)
			}
		}

		// Check if -h or --help are present
		if slices.Contains(arguments, "-h") || slices.Contains(arguments, "--help") {
			ap.UsageFrom(index, parsingState)
			os.Exit(0)
		}

		// Reset all arguments to their default values
		for _, arg := range ap.allArguments {
			arg.ResetDefaultValue()
		}

		// Split between positional arguments and other arguments
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
		missingPositionalArguments := []string{}
		for k, posarg := range ap.PositionalArguments {
			if k < len(potentialPositionalArguments) {
				_, err := posarg.Consume([]string{potentialPositionalArguments[k]})
				if err != nil {
					parsingState.AddErrorMessage(fmt.Sprintf("Error parsing positional argument <%s>: %s", posarg.GetName(), err))
				} else {
					parsingState.ParsedArguments.AddPositionalArgument(&posarg)
				}
			} else {
				missingPositionalArguments = append(missingPositionalArguments, posarg.GetName())
			}
		}
		if len(missingPositionalArguments) != 0 {
			if len(missingPositionalArguments) == 1 {
				parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("Missing %d positional argument: <%s>.", len(missingPositionalArguments), missingPositionalArguments[0]))
			} else {
				errmsg := fmt.Sprintf("Missing %d positional arguments:", len(missingPositionalArguments))
				for _, posarg := range missingPositionalArguments {
					errmsg = errmsg + fmt.Sprintf(" <%s>", posarg)
				}
				errmsg = errmsg + "."
				parsingState.AddErrorMessage(errmsg)
			}
		}
		if len(potentialPositionalArguments) > len(ap.PositionalArguments) {
			leftoverPositionalArguments := potentialPositionalArguments[len(ap.PositionalArguments):]
			if len(leftoverPositionalArguments) == 1 {
				parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("Got %d more positional argument than expected: \"%s\".", len(leftoverPositionalArguments), leftoverPositionalArguments[0]))
			} else {
				errmsg := fmt.Sprintf("Got %d more positional argument than expected: ", len(leftoverPositionalArguments))
				for _, loposarg := range leftoverPositionalArguments {
					errmsg = errmsg + fmt.Sprintf(" \"%s\"", loposarg)
				}
				errmsg = errmsg + "."
				parsingState.AddErrorMessage(errmsg)
			}
		}

		// Parse all other arguments
		// Positions consumed as the value of a recognized flag are tracked so that
		// values which look like flags (e.g. "--port -1") are not reported as unknown.
		consumedAsValue := make(map[int]bool)
		for k, otherarg := range otherArguments {
			if strings.HasPrefix(otherarg, "--") {
				// Long flag name
				if _, exists := ap.longNameToArgument[otherarg]; exists {
					arg := ap.longNameToArgument[otherarg]
					remaining, err := arg.Consume(otherArguments[k:])
					if err != nil {
						parsingState.AddErrorMessage(fmt.Sprintf("Error parsing argument: %s", err))
					} else {
						parsingState.ParsedArguments.AddArgument(&arg)
					}
					for i := k + 1; i < len(otherArguments)-len(remaining); i++ {
						consumedAsValue[i] = true
					}
				} else if !consumedAsValue[k] {
					parsingState.AddErrorMessage(fmt.Sprintf("Unknown argument \"%s\".", otherarg))
				}
			} else if strings.HasPrefix(otherarg, "-") {
				// Short flag name
				if _, exists := ap.shortNameToArgument[otherarg]; exists {
					arg := ap.shortNameToArgument[otherarg]
					remaining, err := arg.Consume(otherArguments[k:])
					if err != nil {
						parsingState.AddErrorMessage(fmt.Sprintf("Error parsing argument: %s", err))
					} else {
						parsingState.ParsedArguments.AddArgument(&arg)
					}
					for i := k + 1; i < len(otherArguments)-len(remaining); i++ {
						consumedAsValue[i] = true
					}
				} else if !consumedAsValue[k] {
					parsingState.AddErrorMessage(fmt.Sprintf("Unknown argument \"%s\".", otherarg))
				}
			}
		}

		// Check if all required arguments have been parsed
		requiredArgumentsMissing := []string{}
		for _, arg := range ap.requiredArguments {
			if !arg.IsPresent() {
				requiredArgumentsMissing = append(requiredArgumentsMissing, arg.GetLongName())
			}
		}
		if len(requiredArgumentsMissing) != 0 {
			if len(requiredArgumentsMissing) == 1 {
				parsingState.AddErrorMessage(fmt.Sprintf("Missing required argument \"%s\"", requiredArgumentsMissing[0]))
			} else {
				parsingState.AddErrorMessage(fmt.Sprintf("Missing required arguments \"%s\"", strings.Join(requiredArgumentsMissing, "\", \"")))
			}
		}

		// Check if all required arguments in groups have been parsed, in a stable order so that
		// the resulting error messages are always reported in the same sequence
		for _, groupName := range ap.sortedGroupNames() {
			group := ap.Groups[groupName]
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
						parsingState.AddErrorMessage(formatGroupErrorMessage(group.Name, fmt.Sprintf("the argument \"%s\" needs to be set.", argumentsMissing[0])))
					} else if len(argumentsMissing) > 1 {
						parsingState.AddErrorMessage(formatGroupErrorMessage(group.Name, fmt.Sprintf("at least one of the arguments \"%s\" needs to be set.", strings.Join(argumentsMissing, "\", \""))))
					}
				} else if len(argumentsPresent) > 1 {
					parsingState.AddErrorMessage(formatGroupErrorMessage(group.Name, fmt.Sprintf("arguments \"%s\" cannot be set together.", strings.Join(argumentsPresent, "\", \""))))
				}
			} else if group.Type == argumentgroup.ARGUMENT_GROUP_TYPE_NOT_REQUIRED_MUTUALLY_EXCLUSIVE {
				// None can be set but if one is set then only one has to be set
				if len(argumentsPresent) > 1 {
					parsingState.AddErrorMessage(formatGroupErrorMessage(group.Name, fmt.Sprintf("arguments \"%s\" cannot be set together.", strings.Join(argumentsPresent, "\", \""))))
				}
			} else if group.Type == argumentgroup.ARGUMENT_GROUP_TYPE_DEPENDENT {
				// If one is set, all need to be set
				if len(argumentsMissing) != 0 {
					if len(argumentsPresent) > 1 {
						parsingState.AddErrorMessage(formatGroupErrorMessage(group.Name, fmt.Sprintf("when arguments \"%s\" are set, \"%s\" need to be set too.", strings.Join(argumentsPresent, "\", \""), strings.Join(argumentsMissing, "\", \""))))
					} else if len(argumentsPresent) == 1 {
						parsingState.AddErrorMessage(formatGroupErrorMessage(group.Name, fmt.Sprintf("when argument \"%s\" is set, \"%s\" need to be set too.", argumentsPresent[0], strings.Join(argumentsMissing, "\", \""))))
					}
				}
			}
		}
	}

	// If there are error messages, print usage and exit
	if len(parsingState.ErrorMessages) != 0 {
		ap.UsageFrom(index, parsingState)
		for _, errmsg := range parsingState.ErrorMessages {
			fmt.Printf("[!] %s\n", errmsg)
		}
		os.Exit(1)
	}
}

// Parse parses the arguments and returns the parsed arguments.
//
// Returns:
// - A map of parsed arguments.
func (ap *ArgumentsParser) Parse() {
	ap.ParsingState.SetRawArguments(os.Args)
	ap.ParseFrom(1, &ap.ParsingState)
}
