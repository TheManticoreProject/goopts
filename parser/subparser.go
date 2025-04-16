package parser

import "strings"

type SubParsers struct {
	// Name is the name of the subparser.
	Name string
	// Value is the value of the subparser.
	Value *string
	// Enabled is a boolean that indicates whether the subparser is enabled.
	Enabled bool
	// CaseInsensitive is a boolean that indicates whether the subparser is case-insensitive.
	CaseInsensitive bool
	// Parsers is a map of subparsers.
	Parsers map[string]*ArgumentsParser
}

// AddSubParser adds a new subparser to the SubParsers.
// It creates a new ArgumentsParser with the specified name and banner.
// If the subparser is case-insensitive, it adds the subparser to the SubParsers with the specified name in lowercase.
// Otherwise, it adds the subparser to the SubParsers with the specified name.
func (sp *SubParsers) AddSubParser(name, banner string) *ArgumentsParser {
	parser_ptr := &ArgumentsParser{
		Banner: banner,
		Options: ArgumentsParserOptions{
			ShowBannerOnHelp: false,
			ShowBannerOnRun:  false,
		},
	}

	if sp.CaseInsensitive {
		sp.Parsers[strings.ToLower(name)] = parser_ptr
	} else {
		sp.Parsers[name] = parser_ptr
	}

	return parser_ptr
}

// GetSubParser returns the subparser with the specified name.
// If the subparser is case-insensitive, it returns the subparser with the specified name in lowercase.
// Otherwise, it returns the subparser with the specified name.
func (sp *SubParsers) GetSubParser(name string) *ArgumentsParser {
	if sp.CaseInsensitive {
		return sp.Parsers[strings.ToLower(name)]
	}
	return sp.Parsers[name]
}

// SetupSubParsing initializes a new subparser with the specified name, value, and case sensitivity.
//
// Parameters:
// - name: The name of the subparser.
// - value: A pointer to a string where the subparser name will be stored.
// - caseInsensitive: A boolean indicating if the subparser name matching should be case insensitive.
//
// Returns:
// - A pointer to the newly created ArgumentsParser instance for the subparser.
func (ap *ArgumentsParser) SetupSubParsing(name string, value *string, caseInsensitive bool) {
	ap.SubParsers = SubParsers{
		Name:            name,
		Value:           value,
		Enabled:         true,
		CaseInsensitive: caseInsensitive,
		Parsers:         make(map[string]*ArgumentsParser),
	}
}

// AddSubParser adds a new subparser to the ArgumentsParser.
//
// Parameters:
// - name: The name of the subparser.
//
// Returns:
// - A pointer to the newly created ArgumentsParser instance for the subparser.
func (ap *ArgumentsParser) AddSubParser(name, banner string) *ArgumentsParser {
	return ap.SubParsers.AddSubParser(name, banner)
}
