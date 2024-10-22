package positionals

import (
	"testing"
)

func TestStringPositionalArgument_Init(t *testing.T) {
	var strValue string
	arg := StringPositionalArgument{}

	arg.Init(&strValue, "name", "Specify a name")

	if arg.GetName() != "name" {
		t.Errorf("Expected name 'name', got '%s'", arg.GetName())
	}
	if arg.GetHelp() != "Specify a name" {
		t.Errorf("Expected help message 'Specify a name', got '%s'", arg.GetHelp())
	}
	if arg.Value != &strValue {
		t.Errorf("Expected Value to point to %v, got %v", strValue, *arg.Value)
	}
}

func TestStringPositionalArgument_GetValue(t *testing.T) {
	var value string = "testValue"
	arg := StringPositionalArgument{
		Value: &value,
	}

	result := arg.GetValue()
	if result != "testValue" {
		t.Errorf("Expected value 'testValue', got %v", result)
	}
}

func TestStringPositionalArgument_IsRequired(t *testing.T) {
	arg := StringPositionalArgument{
		Required: true,
	}
	if !arg.IsRequired() {
		t.Error("Expected IsRequired to return true, got false")
	}

	arg.Required = false
	if arg.IsRequired() {
		t.Error("Expected IsRequired to return false, got true")
	}
}

func TestStringPositionalArgument_Consume(t *testing.T) {
	var strValue string
	arg := StringPositionalArgument{}
	arg.Init(&strValue, "name", "Specify a name")

	// Test when the argument is provided correctly
	remainingArgs, _ := arg.Consume([]string{"John", "otherArg"})

	if strValue != "John" {
		t.Errorf("Expected strValue to be 'John', got '%s'", strValue)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "otherArg" {
		t.Errorf("Expected remaining arguments to be ['otherArg'], got %v", remainingArgs)
	}

	// Reset strValue for the next test
	strValue = ""

	// Test when no arguments are provided
	remainingArgs, _ = arg.Consume([]string{})

	if strValue != "" {
		t.Errorf("Expected strValue to remain '', got '%s'", strValue)
	}
	if len(remainingArgs) != 0 {
		t.Errorf("Expected remaining arguments to be empty, got %v", remainingArgs)
	}
}

func TestStringPositionalArgument_Consume_NoArguments(t *testing.T) {
	var strValue string
	arg := StringPositionalArgument{}
	arg.Init(&strValue, "name", "Specify a name")

	// Test when no arguments are provided
	remainingArgs, _ := arg.Consume([]string{})

	if strValue != "" {
		t.Errorf("Expected strValue to remain '', got '%s'", strValue)
	}
	if len(remainingArgs) != 0 {
		t.Errorf("Expected remaining arguments to be empty, got %v", remainingArgs)
	}
}
