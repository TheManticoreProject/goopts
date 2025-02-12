package arguments

import (
	"testing"
)

func TestIntRangeArgument_Init(t *testing.T) {
	var value int
	arg := IntRangeArgument{}
	arg.Init(&value, "i", "integer", 10, 1, 100, true, "Test integer argument")
	arg.ResetDefaultValue()

	if arg.GetShortName() != "-i" {
		t.Errorf("Expected short name '-i', got '%s'", arg.GetShortName())
	}
	if arg.GetLongName() != "--integer" {
		t.Errorf("Expected long name '--integer', got '%s'", arg.GetLongName())
	}
	if arg.GetDefaultValue() != 10 {
		t.Errorf("Expected default value 10, got %d", arg.GetDefaultValue())
	}
	if !arg.IsRequired() {
		t.Errorf("Expected argument to be required")
	}
	if arg.RangeStart != 1 || arg.RangeStop != 100 {
		t.Errorf("Expected range [1, 100], got [%d, %d]", arg.RangeStart, arg.RangeStop)
	}
}

func TestIntRangeArgument_Consume_Valid(t *testing.T) {
	var value int
	arg := IntRangeArgument{}
	arg.Init(&value, "i", "integer", 10, 1, 100, true, "Test integer argument")
	arg.ResetDefaultValue()

	args := []string{"-i", "50"}
	remainingArgs, err := arg.Consume(args)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if value != 50 {
		t.Errorf("Expected value 50, got %d", value)
	}
	if len(remainingArgs) != 0 {
		t.Errorf("Expected no remaining arguments, got %v", remainingArgs)
	}
}

func TestIntRangeArgument_Consume_OutOfRange(t *testing.T) {
	var value int
	arg := IntRangeArgument{}
	arg.Init(&value, "i", "integer", 10, 1, 100, true, "Test integer argument")
	arg.ResetDefaultValue()

	args := []string{"-i", "150"}
	remainingArgs, err := arg.Consume(args)

	if err == nil {
		t.Errorf("Expected an error, got none")
	}
	if value != arg.DefaultValue {
		t.Errorf("Expected value to remain unchanged, got %d", value)
	}
	if len(remainingArgs) != 2 {
		t.Errorf("Expected remaining arguments to be unchanged, got %v", remainingArgs)
	}
}

func TestIntRangeArgument_Consume_InvalidInteger(t *testing.T) {
	var value int
	arg := IntRangeArgument{}
	arg.Init(&value, "i", "integer", 10, 1, 100, true, "Test integer argument")
	arg.ResetDefaultValue()

	args := []string{"-i", "notAnInteger"}
	remainingArgs, err := arg.Consume(args)

	if err == nil {
		t.Errorf("Expected an error, got none")
	}
	if value != arg.DefaultValue {
		t.Errorf("Expected value to remain unchanged, got %d", value)
	}
	if len(remainingArgs) != 2 {
		t.Errorf("Expected remaining arguments to be unchanged, got %v", remainingArgs)
	}
}
