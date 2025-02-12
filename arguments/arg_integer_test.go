package arguments

import (
	"testing"
)

func TestIntArgument_Init(t *testing.T) {
	var value int

	arg := IntArgument{}
	arg.Init(&value, "n", "number", 10, true, "This is a help message")
	arg.ResetDefaultValue()

	if arg.ShortName != "-n" {
		t.Errorf("Expected ShortName to be '-n', got '%s'", arg.ShortName)
	}
	if arg.LongName != "--number" {
		t.Errorf("Expected LongName to be '--number', got '%s'", arg.LongName)
	}
	if arg.DefaultValue != 10 {
		t.Errorf("Expected DefaultValue to be 10, got '%d'", arg.DefaultValue)
	}
	if arg.Required != true {
		t.Errorf("Expected Required to be true, got '%v'", arg.Required)
	}
	if arg.Help != "This is a help message" {
		t.Errorf("Expected Help to be 'This is a help message', got '%s'", arg.Help)
	}
}

func TestIntArgument_Consume(t *testing.T) {
	value := 0
	arg := IntArgument{
		ShortName: "-n",
		LongName:  "--number",
		Value:     &value,
	}
	arg.ResetDefaultValue()

	arguments := []string{"-n", "42", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if *arg.Value != 42 {
		t.Errorf("Expected Value to be 42, got '%d'", *arg.Value)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "anotherArg" {
		t.Errorf("Expected remaining arguments to be '[anotherArg]', got '%v'", remainingArgs)
	}
}

func TestIntArgument_Consume_InvalidInput(t *testing.T) {
	value := 0
	arg := IntArgument{
		ShortName: "-n",
		LongName:  "--number",
		Value:     &value,
	}
	arg.ResetDefaultValue()

	arguments := []string{"-n", "notAnInt", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if *arg.Value != 0 {
		t.Errorf("Expected Value to remain 0, got '%d'", *arg.Value)
	}
	if len(remainingArgs) != 3 {
		t.Errorf("Expected remaining arguments to be the same as input, got '%v'", remainingArgs)
	}
}

func TestIntArgument_Consume_NoMatch(t *testing.T) {
	value := 0
	arg := IntArgument{
		ShortName: "-n",
		LongName:  "--number",
		Value:     &value,
	}
	arg.ResetDefaultValue()

	arguments := []string{"-x", "42", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if *arg.Value != 0 {
		t.Errorf("Expected Value to remain 0, got '%d'", *arg.Value)
	}
	if len(remainingArgs) != 3 {
		t.Errorf("Expected remaining arguments to be the same as input, got '%v'", remainingArgs)
	}
}

func TestIntArgument_Getters(t *testing.T) {
	value := 10
	arg := IntArgument{
		ShortName:    "-n",
		LongName:     "--number",
		Help:         "This is a help message",
		Value:        &value,
		DefaultValue: 100,
		Required:     true,
	}
	arg.ResetDefaultValue()
	arg.SetValue(value)

	if arg.GetShortName() != "-n" {
		t.Errorf("Expected ShortName to be '-n', got '%s'", arg.GetShortName())
	}
	if arg.GetLongName() != "--number" {
		t.Errorf("Expected LongName to be '--number', got '%s'", arg.GetLongName())
	}
	if arg.GetHelp() != "This is a help message" {
		t.Errorf("Expected Help to be 'This is a help message', got '%s'", arg.GetHelp())
	}
	if arg.GetValue() != 10 {
		t.Errorf("Expected Value to be 10, got '%v'", arg.GetValue())
	}
	if arg.GetDefaultValue() != 100 {
		t.Errorf("Expected DefaultValue to be 100, got '%v'", arg.GetDefaultValue())
	}
	if arg.IsRequired() != true {
		t.Errorf("Expected Required to be true, got '%v'", arg.IsRequired())
	}
}
