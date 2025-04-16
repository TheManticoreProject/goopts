package parser

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/TheManticoreProject/goopts/argumentgroup"
	"github.com/TheManticoreProject/goopts/arguments"
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
	debug := false

	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] ==> Entering\n", index)
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Calling populateMaps()\n", index)
	}
	ap.populateMaps()

	if len(ap.Banner) != 0 && ap.Options.ShowBannerOnRun {
		fmt.Printf("%s\n\n", ap.Banner)
	}

	if debug {
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Check if subparsers are enabled and have parsers\n", index)
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)]   | SubParsers: %+v\n", index, ap.SubParsers)
		fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)]   | SubParsers.Enabled: %+v\n", index, ap.SubParsers.Enabled)
	}
	if ap.SubParsers.Enabled && len(ap.SubParsers.Parsers) != 0 {
		if index < len(parsingState.RawArguments) {
			subparserName := parsingState.RawArguments[index]
			if subparserName == "-h" || subparserName == "--help" {
				ap.UsageFrom(index, parsingState)
				os.Exit(0)
			}
			if debug {
				fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Check if subparser %s exists\n", index, subparserName)
			}
			if asp, exists := ap.SubParsers.Parsers[subparserName]; exists {
				if debug {
					fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Subparser %s exists, parsing it\n", index, subparserName)
				}
				asp.ParseFrom(index+1, parsingState)
				return
			} else {
				if ap.SubParsers.CaseInsensitive {
					parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("No subparser with name \"%s\" was found.", strings.ToLower(subparserName)))
				} else {
					parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("No subparser with name \"%s\" was found.", subparserName))
				}
			}
		} else {
			ap.UsageFrom(index, parsingState)
			os.Exit(1)
		}
	} else {
		// Prepare arguments and split on "=" for `--arg=value`
		if debug {
			fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Prepare arguments and split on \"=\" for `--arg=value`\n", index)
		}
		arguments := []string{}
		for _, arg := range parsingState.RawArguments[index:] {
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
			ap.UsageFrom(index, parsingState)
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
				parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("Missing %d positional argument: <%s>.", len(missingPositionalArguments), missingPositionalArguments[0]))
			} else {
				errmsg := fmt.Sprintf("Missing %d positional arguments:", len(missingPositionalArguments))
				for _, posarg := range missingPositionalArguments {
					errmsg = errmsg + fmt.Sprintf(" <%s>", posarg)
				}
				errmsg = errmsg + "."
				parsingState.ErrorMessages = append(parsingState.ErrorMessages, errmsg)
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
				parsingState.ErrorMessages = append(parsingState.ErrorMessages, errmsg)
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
						parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("Error parsing argument: %s", err))
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
						parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("Error parsing argument: %s", err))
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
				parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("Missing required argument \"%s\"", requiredArgumentsMissing[0]))
			} else {
				parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("Missing required arguments \"%s\"", strings.Join(requiredArgumentsMissing, "\", \"")))
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
						parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("The argument \"%s\" needs to be set.", argumentsMissing[0]))
					} else if len(argumentsMissing) > 1 {
						parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("At least one of the arguments \"%s\" needs to be set.", strings.Join(argumentsMissing, "\", \"")))
					}
				} else if len(argumentsPresent) > 1 {
					parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("Arguments \"%s\" cannot be set together.", strings.Join(argumentsPresent, "\", \"")))
				}
			} else if group.Type == argumentgroup.ARGUMENT_GROUP_TYPE_NOT_REQUIRED_MUTUALLY_EXCLUSIVE {
				// None can be set but if one is set then only one has to be set
				if len(argumentsPresent) > 1 {
					parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("Arguments \"%s\" cannot be set together.", strings.Join(argumentsPresent, "\", \"")))
				}
			} else if group.Type == argumentgroup.ARGUMENT_GROUP_TYPE_DEPENDENT {
				// If one is set, all need to be set
				if len(argumentsMissing) != 0 {
					if len(argumentsPresent) > 1 {
						parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("When arguments \"%s\" are set, \"%s\" need to be set too.", strings.Join(argumentsPresent, "\", \""), strings.Join(argumentsMissing, "\", \"")))
					} else if len(argumentsPresent) == 1 {
						parsingState.ErrorMessages = append(parsingState.ErrorMessages, fmt.Sprintf("When argument \"%s\" is set, \"%s\" need to be set too.", argumentsPresent[0], strings.Join(argumentsMissing, "\", \"")))
					}
				}
			}
		}
	}

	if len(parsingState.ErrorMessages) != 0 {
		if debug {
			fmt.Printf("[debug][ArgumentsParser.ParseFrom(%d)] Error messages found, printing usage and exiting\n", index)
		}
		ap.UsageFrom(index, parsingState)
		for _, errmsg := range parsingState.ErrorMessages {
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
	parsingState := ParsingState{
		RawArguments: os.Args,
	}
	ap.ParseFrom(1, &parsingState)
}
