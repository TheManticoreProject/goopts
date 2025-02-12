package arguments

import (
	"testing"
)

func TestStringArgument_Init(t *testing.T) {
	var value string

	arg := StringArgument{}
	arg.Init(&value, "v", "verbose", "defaultValue", true, "This is a help message")
	arg.ResetDefaultValue()

	if arg.ShortName != "-v" {
		t.Errorf("Expected ShortName to be '-v', got '%s'", arg.ShortName)
	}
	if arg.LongName != "--verbose" {
		t.Errorf("Expected LongName to be '--verbose', got '%s'", arg.LongName)
	}
	if arg.DefaultValue != "defaultValue" {
		t.Errorf("Expected DefaultValue to be 'defaultValue', got '%s'", arg.DefaultValue)
	}
	if arg.Required != true {
		t.Errorf("Expected Required to be true, got '%v'", arg.Required)
	}
	if arg.Help != "This is a help message" {
		t.Errorf("Expected Help to be 'This is a help message', got '%s'", arg.Help)
	}
}

func TestStringArgument_Consume(t *testing.T) {
	value := "initial"
	arg := StringArgument{
		ShortName:    "-v",
		LongName:     "--verbose",
		Value:        &value,
		DefaultValue: value,
	}
	arg.ResetDefaultValue()
	arg.SetValue(value)

	arguments := []string{"-v", "testValue", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if *arg.Value != "testValue" {
		t.Errorf("Expected Value to be 'testValue', got '%s'", *arg.Value)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "anotherArg" {
		t.Errorf("Expected remaining arguments to be '[anotherArg]', got '%v'", remainingArgs)
	}
}

func TestStringArgument_Consume_NoMatch(t *testing.T) {
	value := "initial"
	arg := StringArgument{
		ShortName:    "-v",
		LongName:     "--verbose",
		Value:        &value,
		DefaultValue: value,
	}
	arg.ResetDefaultValue()
	arg.SetValue(value)

	arguments := []string{"-x", "testValue", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if *arg.Value != "initial" {
		t.Errorf("Expected Value to remain 'initial', got '%s'", *arg.Value)
	}
	if len(remainingArgs) != 3 {
		t.Errorf("Expected remaining arguments to be the same as input, got '%v'", remainingArgs)
	}
}

func TestStringArgument_Getters(t *testing.T) {
	value := "value"
	arg := StringArgument{
		ShortName:    "-v",
		LongName:     "--verbose",
		Help:         "This is a help message",
		Value:        &value,
		DefaultValue: "defaultValue",
		Required:     true,
	}
	arg.ResetDefaultValue()

	if arg.GetShortName() != "-v" {
		t.Errorf("Expected ShortName to be '-v', got '%s'", arg.GetShortName())
	}
	if arg.GetLongName() != "--verbose" {
		t.Errorf("Expected LongName to be '--verbose', got '%s'", arg.GetLongName())
	}
	if arg.GetHelp() != "This is a help message" {
		t.Errorf("Expected Help to be 'This is a help message', got '%s'", arg.GetHelp())
	}
	if arg.GetValue() != "value" {
		t.Errorf("Expected Value to be 'value', got '%v'", arg.GetValue())
	}
	if arg.GetDefaultValue() != "defaultValue" {
		t.Errorf("Expected DefaultValue to be 'defaultValue', got '%v'", arg.GetDefaultValue())
	}
	if arg.IsRequired() != true {
		t.Errorf("Expected Required to be true, got '%v'", arg.IsRequired())
	}
}
