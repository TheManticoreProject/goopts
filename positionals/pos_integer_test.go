package positionals

import (
	"testing"
)

func TestIntPositionalArgument_GetValue(t *testing.T) {
	var value int = 42
	arg := IntPositionalArgument{
		Value: &value,
	}

	result := arg.GetValue()
	if result != int(42) {
		t.Errorf("Expected value 42, got %v", result)
	}
}

func TestIntPositionalArgument_IsRequired(t *testing.T) {
	arg := IntPositionalArgument{
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

func TestIntPositionalArgument_Init(t *testing.T) {
	var intValue int
	arg := IntPositionalArgument{}

	arg.Init(&intValue, "number", "Specify a number")

	if arg.GetName() != "number" {
		t.Errorf("Expected name 'number', got '%s'", arg.GetName())
	}
	if arg.GetHelp() != "Specify a number" {
		t.Errorf("Expected help message 'Specify a number', got '%s'", arg.GetHelp())
	}
	if arg.Value != &intValue {
		t.Errorf("Expected Value to point to %v, got %v", intValue, *arg.Value)
	}
}

func TestIntPositionalArgument_Consume(t *testing.T) {
	var intValue int
	arg := IntPositionalArgument{}
	arg.Init(&intValue, "number", "Specify a number")

	// Test when the argument is provided correctly
	remainingArgs, _ := arg.Consume([]string{"42", "otherArg"})

	if intValue != 42 {
		t.Errorf("Expected intValue to be 42, got %d", intValue)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "otherArg" {
		t.Errorf("Expected remaining arguments to be ['otherArg'], got %v", remainingArgs)
	}

	// Reset intValue for the next test
	intValue = 0

	// Test when the long argument is provided
	remainingArgs, _ = arg.Consume([]string{"123", "anotherArg"})

	if intValue != 123 {
		t.Errorf("Expected intValue to be 123, got %d", intValue)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "anotherArg" {
		t.Errorf("Expected remaining arguments to be ['anotherArg'], got %v", remainingArgs)
	}

	// Reset intValue for the next test
	intValue = 0

	// Test when the argument is not provided
	remainingArgs, _ = arg.Consume([]string{"someOtherArg"})

	if intValue != 0 {
		t.Errorf("Expected intValue to remain 0, got %d", intValue)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "someOtherArg" {
		t.Errorf("Expected remaining arguments to be ['someOtherArg'], got %v", remainingArgs)
	}
}

func TestIntPositionalArgument_Consume_InvalidNumber(t *testing.T) {
	var intValue int
	arg := IntPositionalArgument{}
	arg.Init(&intValue, "number", "Specify a number")

	// Test when an invalid number is provided
	remainingArgs, _ := arg.Consume([]string{"-n", "invalid", "otherArg"})

	if intValue != 0 {
		t.Errorf("Expected intValue to remain 0, got %d", intValue)
	}
	if len(remainingArgs) != 3 || remainingArgs[0] != "-n" || remainingArgs[1] != "invalid" || remainingArgs[2] != "otherArg" {
		t.Errorf("Expected remaining arguments to be ['-n', 'invalid', 'otherArg'], got %v", remainingArgs)
	}
}

func TestIntPositionalArgument_Consume_EmptyArguments(t *testing.T) {
	var intValue int
	arg := IntPositionalArgument{}
	arg.Init(&intValue, "number", "Specify a number")

	remainingArgs, _ := arg.Consume([]string{})

	if intValue != 0 {
		t.Errorf("Expected intValue to remain 0, got %d", intValue)
	}
	if len(remainingArgs) != 0 {
		t.Errorf("Expected remaining arguments to be empty, got %v", remainingArgs)
	}
}
