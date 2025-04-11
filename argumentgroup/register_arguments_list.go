package argumentgroup

import (
	"github.com/TheManticoreProject/goopts/arguments"
)

// NewListOfIntsArgument registers a new list of integers argument with the argument parser.
//
// Parameters:
// - ptr: A pointer to the slice of int variables where the argument values will be stored.
// - shortName: The short name (single character) of the argument, prefixed with a dash (e.g., "-n").
// - longName: The long name of the argument, prefixed with two dashes (e.g., "--numbers").
// - defaultValue: The default list of integers if the argument is not provided by the user.
// - required: A boolean indicating if the argument is mandatory.
// - help: A description of the argument, which will be displayed in the help message.
//
// The function creates a new ListOfIntsArgument with the provided parameters and adds it to the argument group.
func (ag *ArgumentGroup) NewListOfIntsArgument(ptr *[]int, shortName, longName string, defaultValue []int, required bool, help string) error {
	arg := &arguments.ListOfIntsArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, required, help)
	err := ag.Register(arg)
	return err
}

// NewListOfStringsArgument registers a new list of strings argument with the argument parser.
//
// Parameters:
// - ptr: A pointer to the slice of string variables where the argument values will be stored.
// - shortName: The short name (single character) of the argument, prefixed with a dash (e.g., "-s").
// - longName: The long name of the argument, prefixed with two dashes (e.g., "--strings").
// - defaultValue: The default list of strings if the argument is not provided by the user.
// - required: A boolean indicating if the argument is mandatory.
// - help: A description of the argument, which will be displayed in the help message.
//
// The function creates a new ListOfStringsArgument with the provided parameters and adds it to the argument group.
func (ag *ArgumentGroup) NewListOfStringsArgument(ptr *[]string, shortName, longName string, defaultValue []string, required bool, help string) error {
	arg := &arguments.ListOfStringsArgument{}
	arg.Init(ptr, shortName, longName, defaultValue, required, help)
	err := ag.Register(arg)
	return err
}

// NewMapOfHttpHeadersArgument registers a new list of HTTP headers argument with the argument parser.
//
// Parameters:
// - ptr: A pointer to the slice of string variables where the argument values will be stored.
// - shortName: The short name (single character) of the argument, prefixed with a dash (e.g., "-H").
// - longName: The long name of the argument, prefixed with two dashes (e.g., "--header").
// - required: A boolean indicating if the argument is mandatory.
// - help: A description of the argument, which will be displayed in the help message.
//
// The function creates a new MapOfHttpHeadersArgument with the provided parameters and adds it to the argument group.
// Each header should be passed in the format "Key: Value".
func (ag *ArgumentGroup) NewMapOfHttpHeadersArgument(ptr *map[string]string, shortName, longName string, defaultValue map[string]string, required bool, help string) error {
	arg := &arguments.MapOfHttpHeadersArgument{}
	// Default value is an empty list since headers may vary widely
	arg.Init(ptr, shortName, longName, defaultValue, required, help)
	err := ag.Register(arg)
	return err
}
