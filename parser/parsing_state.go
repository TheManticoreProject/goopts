package parser

import "github.com/TheManticoreProject/goopts/arguments"

type ParsingState struct {
	RawArguments []string

	ErrorMessages []string

	ParsedArguments map[string]*arguments.Argument
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

// GetParsedArguments returns the parsed arguments from the parsing state.
//
// Returns:
// - A map of parsed arguments.
func (ps *ParsingState) GetParsedArguments() map[string]*arguments.Argument {
	return ps.ParsedArguments
}

// SetParsedArguments sets the parsed arguments in the parsing state.
//
// Parameters:
// - parsedArguments: A map of parsed arguments.
func (ps *ParsingState) SetParsedArguments(parsedArguments map[string]*arguments.Argument) {
	ps.ParsedArguments = parsedArguments
}
