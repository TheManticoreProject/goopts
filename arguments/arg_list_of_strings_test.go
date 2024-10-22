package arguments

import (
	"fmt"
	"slices"
	"testing"
)

func TestListOfStringsArgument_Init(t *testing.T) {
	var values []string

	arg := ListOfStringsArgument{}
	arg.Init(&values, "s", "strings", []string{"default1", "default2"}, true, "List of strings")

	if arg.ShortName != "-s" {
		t.Errorf("Expected ShortName to be '-s', got '%s'", arg.ShortName)
	}
	if arg.LongName != "--strings" {
		t.Errorf("Expected LongName to be '--strings', got '%s'", arg.LongName)
	}
	if !arg.Required {
		t.Errorf("Expected Required to be true, got '%v'", arg.Required)
	}
	if arg.Help != "List of strings" {
		t.Errorf("Expected Help to be 'List of strings', got '%s'", arg.Help)
	}
	if len(arg.DefaultValue) != 2 || arg.DefaultValue[0] != "default1" || arg.DefaultValue[1] != "default2" {
		t.Errorf("Expected DefaultValue to be [default1, default2], got '%v'", arg.DefaultValue)
	}
	if !slices.Equal(*arg.Value, arg.DefaultValue) {
		t.Errorf("Expected Value to be equal to DefaultValue, got '%v'", *arg.Value)
	}
}

func TestListOfStringsArgument_Consume(t *testing.T) {
	values := []string{}
	arg := ListOfStringsArgument{
		ShortName: "-s",
		LongName:  "--strings",
		Value:     &values,
	}

	arguments := []string{"-s", "example", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if len(*arg.Value) != 1 || (*arg.Value)[0] != "example" {
		t.Errorf("Expected Value to contain ['example'], got '%v'", *arg.Value)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "anotherArg" {
		t.Errorf("Expected remaining arguments to be '[anotherArg]', got '%v'", remainingArgs)
	}
}

func TestListOfStringsArgument_Consume_Multiple(t *testing.T) {
	values := []string{}
	arg := ListOfStringsArgument{
		ShortName: "-s",
		LongName:  "--strings",
		Value:     &values,
	}

	arguments := []string{"-s", "example1", "-s", "example2", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)
	remainingArgs, _ = arg.Consume(remainingArgs)

	if len(*arg.Value) != 2 || (*arg.Value)[0] != "example1" || (*arg.Value)[1] != "example2" {
		t.Errorf("Expected Value to contain ['example1', 'example2'], got '%v'", *arg.Value)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "anotherArg" {
		t.Errorf("Expected remaining arguments to be '[anotherArg]', got '%v'", remainingArgs)
	}
}

func TestListOfStringsArgument_Consume_NoMatch(t *testing.T) {
	values := []string{}
	arg := ListOfStringsArgument{
		ShortName: "-s",
		LongName:  "--strings",
		Value:     &values,
	}

	arguments := []string{"-x", "example", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if len(*arg.Value) != 0 {
		t.Errorf("Expected Value to remain empty, got '%v'", *arg.Value)
	}
	if len(remainingArgs) != 3 {
		t.Errorf("Expected remaining arguments to be the same as input, got '%v'", remainingArgs)
	}
}

func TestListOfStringsArgument_Getters(t *testing.T) {
	values := []string{"item1", "item2"}
	arg := ListOfStringsArgument{
		ShortName:    "-s",
		LongName:     "--strings",
		Help:         "List of strings",
		Value:        &values,
		DefaultValue: []string{"default1", "default2"},
		Required:     true,
	}

	if arg.GetShortName() != "-s" {
		t.Errorf("Expected ShortName to be '-s', got '%s'", arg.GetShortName())
	}
	if arg.GetLongName() != "--strings" {
		t.Errorf("Expected LongName to be '--strings', got '%s'", arg.GetLongName())
	}
	if arg.GetHelp() != "List of strings" {
		t.Errorf("Expected Help to be 'List of strings', got '%s'", arg.GetHelp())
	}
	if val, ok := arg.GetValue().([]string); ok {
		if !slices.Equal(val, values) {
			fmt.Printf("Expected Value to be '%v', got '%v'\n", values, val)
		}
	} else {
		fmt.Printf("Expected GetValue to return a []string, but got a different type\n")
	}
	if val, ok := arg.GetDefaultValue().([]string); ok {
		if !slices.Equal(val, []string{"default1", "default2"}) {
			t.Errorf("Expected DefaultValue to be [default1, default2], got '%v'", arg.GetDefaultValue())
		}
	} else {
		fmt.Printf("Expected GetValue to return a []string, but got a different type\n")
	}
	if !arg.IsRequired() {
		t.Errorf("Expected Required to be true, got '%v'", arg.IsRequired())
	}
}
