package positionals

import (
	"testing"
)

func TestBoolPositionalArgument_Init(t *testing.T) {
	var boolValue bool
	arg := BoolPositionalArgument{}

	arg.Init(&boolValue, "force", "Force operation")

	if arg.GetName() != "force" {
		t.Errorf("Expected long name 'force', got '%s'", arg.GetName())
	}
	if arg.GetHelp() != "Force operation" {
		t.Errorf("Expected help message 'Force operation', got '%s'", arg.GetHelp())
	}
	if arg.Value != &boolValue {
		t.Errorf("Expected Value to point to %v, got %v", boolValue, *arg.Value)
	}
}

func TestBoolPositionalArgument_GetValue(t *testing.T) {
	var value bool = true
	arg := BoolPositionalArgument{
		Value: &value,
	}

	result := arg.GetValue()
	if result != true {
		t.Errorf("Expected value true, got %v", result)
	}
}

func TestBoolPositionalArgument_IsRequired(t *testing.T) {
	arg := BoolPositionalArgument{
		Required: true,
	}
	if !arg.IsRequired() {
		t.Error("Expected IsRequired to return true, got false")
	}

	arg.Required = false
	if !arg.IsRequired() {
		t.Error("Expected IsRequired to return true, got false")
	}
}

func TestBoolPositionalArgument_Consume(t *testing.T) {
	var boolValue bool
	arg := BoolPositionalArgument{}
	arg.Init(&boolValue, "force", "Force operation")

	// Test when the argument is provided
	remainingArgs, _ := arg.Consume([]string{"true", "otherArg"})

	if !boolValue {
		t.Errorf("Expected boolValue to be true, got false")
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "otherArg" {
		t.Errorf("Expected remaining arguments to be ['otherArg'], got %v", remainingArgs)
	}
}

func TestBoolPositionalArgument_Consume_EmptyArguments(t *testing.T) {
	var boolValue bool
	arg := BoolPositionalArgument{}
	arg.Init(&boolValue, "force", "Force operation")

	remainingArgs, _ := arg.Consume([]string{})

	if boolValue {
		t.Errorf("Expected boolValue to be false, got true")
	}
	if len(remainingArgs) != 0 {
		t.Errorf("Expected remaining arguments to be empty, got %v", remainingArgs)
	}
}
