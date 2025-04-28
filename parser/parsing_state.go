package parser

import (
	"github.com/TheManticoreProject/goopts/arguments"
	"github.com/TheManticoreProject/goopts/positionals"
)

type ParsingState struct {
	RawArguments    []string
	ErrorMessages   []string
	ParsedArguments ParsedArguments
}

type ParsedArguments struct {
	PositionalArguments map[string]*positionals.PositionalArgument
	LongNameToArgument  map[string]*arguments.Argument
	ShortNameToArgument map[string]*arguments.Argument
}

// SetRawArguments sets the raw arguments in the parsing state.
//
// Parameters:
// - rawArguments: The raw arguments to set.
func (ps *ParsingState) SetRawArguments(rawArguments []string) {
	ps.RawArguments = rawArguments
}

// AddErrorMessage adds an error message to the parsing state.
//
// Parameters:
// - message: The error message to add.
func (ps *ParsingState) AddErrorMessage(message string) {
	ps.ErrorMessages = append(ps.ErrorMessages, message)
}

// ClearErrorMessages clears the error messages from the parsing state.
//
// This method resets the error messages slice to an empty state, effectively clearing any previously added error messages.
func (ps *ParsingState) ClearErrorMessages() {
	ps.ErrorMessages = []string{}
}

// GetErrorMessages returns the error messages from the parsing state.
//
// Returns:
// - A slice of error messages.
func (ps *ParsingState) GetErrorMessages() []string {
	return ps.ErrorMessages
}

// AddArgument adds an argument to the parsing state.
//
// Parameters:
// - argument: The argument to add.
func (pa *ParsedArguments) AddArgument(argument *arguments.Argument) {
	longName := (*argument).GetLongName()
	shortName := (*argument).GetShortName()

	if len(longName) != 0 {
		pa.LongNameToArgument[longName] = argument
	}
	if len(shortName) != 0 {
		pa.ShortNameToArgument[shortName] = argument
	}
}

// AddPositionalArgument adds a positional argument to the parsing state.
//
// Parameters:
// - argument: The positional argument to add.
func (pa *ParsedArguments) AddPositionalArgument(argument *positionals.PositionalArgument) {
	name := (*argument).GetName()
	pa.PositionalArguments[name] = argument
}
