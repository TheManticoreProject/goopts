package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// StripLeftDashes removes all leading dash characters ('-') from the input string.
// This function is useful for sanitizing command-line arguments or any other input
// that may contain leading dashes, such as flags or options.
//
// Parameters:
//
//	s (string): The input string from which to strip leading dashes.
//
// Returns:
//
//	string: A new string with all leading dashes removed.
func StripLeftDashes(s string) string {
	for len(s) > 0 && s[0] == '-' {
		s = s[1:]
	}
	return s
}

// GenerateLongAndShortNames generates the long and short names for an argument.
//
// Parameters:
//   - longName (string): The long name of the argument.
//   - shortName (string): The short name of the argument.
//
// Returns:
//   - string: The long name of the argument.
//   - string: The short name of the argument.
func GenerateLongAndShortNames(longName, shortName string) (string, string) {
	if len(shortName) == 0 {
		shortName = ""
	} else {
		shortName = "-" + StripLeftDashes(shortName)
	}

	if len(longName) == 0 {
		longName = ""
	} else {
		longName = "--" + StripLeftDashes(longName)
	}

	return longName, shortName
}

// ListOfStrings converts a slice of strings into a formatted string representation.
// Each element of the slice is enclosed in double quotes and separated by commas.
// The entire output is enclosed in square brackets, mimicking a JSON-style array.
//
// Parameters:
//
//	l ([]string): The slice of strings to be formatted.
//
// Returns:
//
//	string: A formatted string representation of the slice, with each element quoted and separated by commas.
//	        Example: ["item1", "item2", "item3"]
func ListOfStrings(l []string) string {
	outputString := "["

	for k, a := range l {
		if k == (len(l) - 1) {
			outputString = fmt.Sprintf("%s\"%s\"", outputString, a)
		} else {
			outputString = fmt.Sprintf("%s\"%s\", ", outputString, a)
		}
	}
	outputString = fmt.Sprintf("%s]", outputString)

	return outputString
}

// StringToInt converts a string representation of an integer in various formats
// (decimal, hexadecimal, octal, or binary) to an integer.
// The function recognizes the following prefixes:
// - "0x" or "0X" for hexadecimal
// - "0o" or "0O" for octal
// - "0b" or "0B" for binary
// If no prefix is provided, the string is treated as a decimal.
//
// Parameters:
// - value: The string to convert to an integer.
//
// Returns:
// - The converted integer value.
// - An error if the conversion fails.
func StringToInt(value string) (int, error) {
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		// Hexadecimal
		parsed, err := strconv.ParseInt(value[2:], 16, 64)
		return int(parsed), err
	} else if strings.HasPrefix(value, "0o") || strings.HasPrefix(value, "0O") {
		// Octal
		parsed, err := strconv.ParseInt(value[2:], 8, 64)
		return int(parsed), err
	} else if strings.HasPrefix(value, "0b") || strings.HasPrefix(value, "0B") {
		// Binary
		parsed, err := strconv.ParseInt(value[2:], 2, 64)
		return int(parsed), err
	} else {
		// Decimal
		return strconv.Atoi(value)
	}
}
