package parser

import "github.com/TheManticoreProject/goopts/arguments"

type ParsingState struct {
	RawArguments []string

	ErrorMessages []string

	ParsedArguments map[string]*arguments.Argument
}
