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

	if ap.Groups == nil {
		ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)
	}

	// Register arguments from the default group
	if ap.Groups[""] != nil {
		for _, arg := range ap.Groups[""].Arguments {
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
		ap.Groups[""] = &ag
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

	if ap.ParsingState.ParsedArguments.PositionalArguments == nil {
		ap.ParsingState.ParsedArguments.PositionalArguments = make(map[string]*positionals.PositionalArgument)
	}
	if ap.ParsingState.ParsedArguments.LongNameToArgument == nil {
		ap.ParsingState.ParsedArguments.LongNameToArgument = make(map[string]*arguments.Argument)
	}
	if ap.ParsingState.ParsedArguments.ShortNameToArgument == nil {
		ap.ParsingState.ParsedArguments.ShortNameToArgument = make(map[string]*arguments.Argument)
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
func (ap *ArgumentsParser) ParseFrom(index int, parsingState *ParsingState) {
	ap.populateMaps()

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
			if asp, exists := ap.SubParsers.Parsers[subparserName]; exists {
				// Set the subparser name value to the pointer
				*(ap.SubParsers.Value) = subparserName
				asp.ParseFrom(index+1, parsingState)
				return
			} else {
				if ap.SubParsers.CaseInsensitive {
					parsingState.AddErrorMessage(fmt.Sprintf("No subparser with name \"%s\" was found.", strings.ToLower(subparserName)))
				} else {
					parsingState.AddErrorMessage(fmt.Sprintf("No subparser with name \"%s\" was found.", subparserName))
				}
			}
		} else {
			ap.UsageFrom(index, parsingState)
			os.Exit(1)
		}
	} else {
		// Prepare arguments and split on "=" for `--arg=value`
		arguments := []string{}
		for _, arg := range parsingState.RawArguments[index:] {
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
				posarg.Consume([]string{potentialPositionalArguments[k]})
				parsingState.ParsedArguments.AddPositionalArgument(&posarg)
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
		for k, otherarg := range otherArguments {
			if strings.HasPrefix(otherarg, "--") {
				// Long flag name
				if _, exists := ap.longNameToArgument[otherarg]; exists {
					arg := ap.longNameToArgument[otherarg]
					_, err := arg.Consume(otherArguments[k:])
					if err != nil {
						parsingState.AddErrorMessage(fmt.Sprintf("Error parsing argument: %s", err))
					} else {
						parsingState.ParsedArguments.AddArgument(&arg)
					}
				}
			} else if strings.HasPrefix(otherarg, "-") {
				// Short flag name
				if _, exists := ap.shortNameToArgument[otherarg]; exists {
					arg := ap.shortNameToArgument[otherarg]
					_, err := arg.Consume(otherArguments[k:])
					if err != nil {
						parsingState.AddErrorMessage(fmt.Sprintf("Error parsing argument: %s", err))
					} else {
						parsingState.ParsedArguments.AddArgument(&arg)
					}
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

		// Check if all required arguments in groups have been parsed
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
						parsingState.AddErrorMessage(fmt.Sprintf("The argument \"%s\" needs to be set.", argumentsMissing[0]))
					} else if len(argumentsMissing) > 1 {
						parsingState.AddErrorMessage(fmt.Sprintf("At least one of the arguments \"%s\" needs to be set.", strings.Join(argumentsMissing, "\", \"")))
					}
				} else if len(argumentsPresent) > 1 {
					parsingState.AddErrorMessage(fmt.Sprintf("Arguments \"%s\" cannot be set together.", strings.Join(argumentsPresent, "\", \"")))
				}
			} else if group.Type == argumentgroup.ARGUMENT_GROUP_TYPE_NOT_REQUIRED_MUTUALLY_EXCLUSIVE {
				// None can be set but if one is set then only one has to be set
				if len(argumentsPresent) > 1 {
					parsingState.AddErrorMessage(fmt.Sprintf("Arguments \"%s\" cannot be set together.", strings.Join(argumentsPresent, "\", \"")))
				}
			} else if group.Type == argumentgroup.ARGUMENT_GROUP_TYPE_DEPENDENT {
				// If one is set, all need to be set
				if len(argumentsMissing) != 0 {
					if len(argumentsPresent) > 1 {
						parsingState.AddErrorMessage(fmt.Sprintf("When arguments \"%s\" are set, \"%s\" need to be set too.", strings.Join(argumentsPresent, "\", \""), strings.Join(argumentsMissing, "\", \"")))
					} else if len(argumentsPresent) == 1 {
						parsingState.AddErrorMessage(fmt.Sprintf("When argument \"%s\" is set, \"%s\" need to be set too.", argumentsPresent[0], strings.Join(argumentsMissing, "\", \"")))
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
