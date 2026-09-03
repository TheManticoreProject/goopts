package parser

import (
	"sort"
	"strings"
)

type SubParsers struct {
	// Name is the name of the subparser.
	Name string
	// Value is the value of the subparser.
	Value *string
	// Enabled is a boolean that indicates whether the subparser is enabled.
	Enabled bool
	// CaseInsensitive is a boolean that indicates whether the subparser is case-insensitive.
	//
	// The zero value is case-sensitive matching. Names are matched at lookup time and never
	// normalized when they are registered, so this can be flipped at any point of the setup,
	// before or after the subparsers themselves are registered.
	CaseInsensitive bool
	// Parsers is a map of subparsers, keyed by the name they were registered with.
	Parsers map[string]*ArgumentsParser
}

// resolve returns the subparser matching name, along with the name it was registered with.
//
// An exact match always wins. When matching is case-insensitive and no exact match exists, the
// name is compared against every registered name without regard to case. Several registered names
// can differ only by case, so the candidates are sorted and the first one is returned to make the
// dispatch deterministic instead of following map iteration order.
func (sp *SubParsers) resolve(name string) (string, *ArgumentsParser, bool) {
	if parser_ptr, exists := sp.Parsers[name]; exists {
		return name, parser_ptr, true
	}

	if !sp.CaseInsensitive {
		return "", nil, false
	}

	candidates := []string{}
	for registeredName := range sp.Parsers {
		if strings.EqualFold(registeredName, name) {
			candidates = append(candidates, registeredName)
		}
	}
	if len(candidates) == 0 {
		return "", nil, false
	}
	sort.Strings(candidates)

	return candidates[0], sp.Parsers[candidates[0]], true
}

// AddSubParser adds a new subparser to the SubParsers.
// It creates a new ArgumentsParser with the specified name and banner.
//
// The subparser is registered under the name it is given, whatever the case sensitivity setting
// is: names are compared when they are looked up, not rewritten when they are registered, so the
// name stays available for the usage output and for the value SetupSubParsing binds.
//
// The new parser inherits the case sensitivity of the parser it is registered on, and can override
// it with SetSubparserNameCaseSensitive.
//
// The map of subparsers is created if needed and subparsing is enabled, so this can be called
// without SetupSubParsing having run first, which is the case for every subparser nested under
// another one. SetupSubParsing is still what binds a pointer to receive the selected subparser
// name; without it the name is matched and dispatched but not stored anywhere.
func (sp *SubParsers) AddSubParser(name, banner string) *ArgumentsParser {
	parser_ptr := &ArgumentsParser{
		Banner: banner,
		Options: ArgumentsParserOptions{
			ShowBannerOnHelp: false,
			ShowBannerOnRun:  false,
		},
		SubParsers: SubParsers{
			CaseInsensitive: sp.CaseInsensitive,
		},
	}

	// Only SetupSubParsing allocates the map, and the parsers created here never went through it,
	// so a nested subparser would otherwise write to a nil map. Registering a subparser is also
	// intent to dispatch on one, which is what Enabled means.
	if sp.Parsers == nil {
		sp.Parsers = make(map[string]*ArgumentsParser)
	}
	sp.Enabled = true

	sp.Parsers[name] = parser_ptr

	return parser_ptr
}

// GetSubParser returns the subparser with the specified name, or nil when no subparser matches it.
// The name is matched exactly, unless case-insensitive matching is enabled.
func (sp *SubParsers) GetSubParser(name string) *ArgumentsParser {
	_, parser_ptr, exists := sp.resolve(name)
	if !exists {
		return nil
	}
	return parser_ptr
}

// SetupSubParsing initializes subparsing on the parser with the specified name, value, and case
// sensitivity. Subparsers already registered on the parser are kept.
//
// Parameters:
// - name: The name of the subparser.
// - value: A pointer to a string where the name of the selected subparser will be stored, as it
// was registered with AddSubParser.
// - caseInsensitive: A boolean indicating if the subparser name matching should be case
// insensitive. This is the same setting as SetSubparserNameCaseSensitive, with the opposite
// polarity, and it likewise applies to the subparsers already registered on the parser.
func (ap *ArgumentsParser) SetupSubParsing(name string, value *string, caseInsensitive bool) {
	ap.SubParsers.Name = name
	ap.SubParsers.Value = value
	ap.SubParsers.Enabled = true
	if ap.SubParsers.Parsers == nil {
		ap.SubParsers.Parsers = make(map[string]*ArgumentsParser)
	}
	ap.SetSubparserNameCaseSensitive(!caseInsensitive)
}

// SetSubparserNameCaseSensitive sets whether the subparser names of this parser must be given with
// the case they were registered with. Subparser names are case-sensitive by default.
//
// It can be called before or after AddSubParser, because names are matched when they are looked up
// and never normalized when they are registered. It applies to every subparser already registered
// below this parser, and subparsers registered afterwards inherit it, so a single call on the root
// parser configures the whole tree. A subparser overrides the setting for its own subparsers by
// calling this on itself after it has been registered.
//
// Parameters:
// - caseSensitive: A boolean indicating whether subparser names are matched case-sensitively.
func (ap *ArgumentsParser) SetSubparserNameCaseSensitive(caseSensitive bool) {
	ap.SubParsers.CaseInsensitive = !caseSensitive

	for _, subparser := range ap.SubParsers.Parsers {
		subparser.SetSubparserNameCaseSensitive(caseSensitive)
	}
}

// SubparserNameIsCaseSensitive returns whether the subparser names of this parser are matched
// case-sensitively.
func (ap *ArgumentsParser) SubparserNameIsCaseSensitive() bool {
	return !ap.SubParsers.CaseInsensitive
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
