package arguments

import (
	"testing"
)

func TestBoolArgument_Init(t *testing.T) {
	var value bool

	arg := BoolArgument{}
	arg.Init(&value, "v", "verbose", false, "Enable verbose mode")
	arg.ResetDefaultValue()

	if arg.ShortName != "-v" {
		t.Errorf("Expected ShortName to be '-v', got '%s'", arg.ShortName)
	}
	if arg.LongName != "--verbose" {
		t.Errorf("Expected LongName to be '--verbose', got '%s'", arg.LongName)
	}
	if arg.DefaultValue != false {
		t.Errorf("Expected DefaultValue to be false, got '%v'", arg.DefaultValue)
	}
	if arg.Help != "Enable verbose mode" {
		t.Errorf("Expected Help to be 'Enable verbose mode', got '%s'", arg.Help)
	}
	if *arg.Value != false {
		t.Errorf("Expected Value to be initialized to false, got '%v'", *arg.Value)
	}
}

func TestBoolArgument_Consume(t *testing.T) {
	value := false
	arg := BoolArgument{
		ShortName:    "-v",
		LongName:     "--verbose",
		Value:        &value,
		DefaultValue: value,
	}
	arg.ResetDefaultValue()
	arg.SetValue(value)

	arguments := []string{"-v", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if *arg.Value != true {
		t.Errorf("Expected Value to be true, got '%v'", *arg.Value)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "anotherArg" {
		t.Errorf("Expected remaining arguments to be '[anotherArg]', got '%v'", remainingArgs)
	}
}

func TestBoolArgument_Consume_NoMatch(t *testing.T) {
	value := false
	arg := BoolArgument{
		ShortName:    "-v",
		LongName:     "--verbose",
		Value:        &value,
		DefaultValue: value,
	}
	arg.ResetDefaultValue()
	arg.SetValue(value)

	arguments := []string{"-x", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if *arg.Value != false {
		t.Errorf("Expected Value to remain false, got '%v'", *arg.Value)
	}
	if len(remainingArgs) != 2 {
		t.Errorf("Expected remaining arguments to be the same as input, got '%v'", remainingArgs)
	}
}

func TestBoolArgument_Getters(t *testing.T) {
	value := true
	arg := BoolArgument{
		ShortName:    "-v",
		LongName:     "--verbose",
		Help:         "Enable verbose mode",
		Value:        &value,
		DefaultValue: false,
	}
	arg.ResetDefaultValue()
	arg.SetValue(value)

	if arg.GetShortName() != "-v" {
		t.Errorf("Expected ShortName to be '-v', got '%s'", arg.GetShortName())
	}
	if arg.GetLongName() != "--verbose" {
		t.Errorf("Expected LongName to be '--verbose', got '%s'", arg.GetLongName())
	}
	if arg.GetHelp() != "Enable verbose mode" {
		t.Errorf("Expected Help to be 'Enable verbose mode', got '%s'", arg.GetHelp())
	}
	if arg.GetValue() != true {
		t.Errorf("Expected Value to be true, got '%v'", arg.GetValue())
	}
	if arg.GetDefaultValue() != false {
		t.Errorf("Expected DefaultValue to be false, got '%v'", arg.GetDefaultValue())
	}
	if arg.IsRequired() != false {
		t.Errorf("Expected Required to be false, got '%v'", arg.IsRequired())
	}
}
