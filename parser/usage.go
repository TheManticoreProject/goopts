package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/TheManticoreProject/goopts/arguments"
)

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
func (ap *ArgumentsParser) UsageFrom(index int, parsingState *ParsingState) {
	// Create usage string
	usage := "Usage: " + filepath.Base(parsingState.RawArguments[0])
	if index != 0 {
		for k := range index - 1 {
			usage += " " + parsingState.RawArguments[k+1]
		}
	}

	// Add subparsers
	if ap.SubParsers.Enabled {
		names := []string{}
		for name := range ap.SubParsers.Parsers {
			names = append(names, name)
		}
		sort.Strings(names)
		usage += " <" + strings.Join(names, "|") + ">"

		// Compute the maximum length of the subparser names
		maxLen := 0
		for _, name := range names {
			if len(name) > maxLen {
				maxLen = len(name)
			}
		}
		// Print the subparsers
		usage += "\n\n"
		fmtString := fmt.Sprintf("   %%-%ds  %%s\n", maxLen)
		for _, name := range names {
			usage += fmt.Sprintf(fmtString, name, ap.SubParsers.Parsers[name].Banner)
		}

	} else {
		// This is the usage line ============================================================
		// Add positional arguments
		for _, posarg := range ap.PositionalArguments {
			usage += fmt.Sprintf(" <%s>", posarg.GetName())
		}
		// Append default group arguments
		for _, argument := range ap.Groups[""].Arguments {
			output := generateArgumentForUsageLine(argument)
			if len(output) != 0 {
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
			// This is the default group, we don't want to print it in this block
			if groupname == "" {
				continue
			}
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
		usage += "\n\n"

		// This is the detailled help for each group ============================================================
		for _, groupname := range groupNames {
			group := ap.Groups[groupname]

			// This is the default group, we want to print differently
			if groupname == "" {
				fmtString := fmt.Sprintf("  %s %%s\n", computePaddingFormat(ap.Groups[""].Arguments))
				for _, argument := range group.Arguments {
					usage += generateArgumentLineInHelp(argument, fmtString)
				}
			} else {
				fmtString := fmt.Sprintf("    %s %%s\n", computePaddingFormat(group.Arguments))
				usage += fmt.Sprintf("\n  %s:\n", group.Name)
				for _, argument := range group.Arguments {
					usage += generateArgumentLineInHelp(argument, fmtString)
				}
			}

		}
	}

	// Finally print the usage
	fmt.Printf("%s\n", usage)
}

// Usage prints the usage information for the command-line arguments.
func (ap *ArgumentsParser) Usage() {
	// We start printing from index 0 because the first argument (0)
	// is the program name and we need it in the usage
	parsingState := &ParsingState{
		RawArguments: os.Args,
	}
	ap.UsageFrom(0, parsingState)
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
func generateArgumentLineInHelp(arg arguments.Argument, fmtString string) string {
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

	return fmt.Sprintf(fmtString, flags_string, help)
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
