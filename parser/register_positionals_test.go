package parser

import (
	"testing"

	"github.com/TheManticoreProject/goopts/positionals"
)

func TestRegisterStringPositionalArgument(t *testing.T) {
	ap := &ArgumentsParser{}
	ap.PositionalArguments = make([]positionals.PositionalArgument, 0)

	var positionalArg string
	err := ap.NewStringPositionalArgument(&positionalArg, "username", "The username argument")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(ap.PositionalArguments) != 1 {
		t.Errorf("Expected 1 positional argument, got %d", len(ap.PositionalArguments))
	}
	arg := ap.PositionalArguments[0]
	if arg.GetName() != "username" {
		t.Errorf("Expected positional argument name 'username', got %s", arg.GetName())
	}

	// Test duplicate registration
	err = ap.NewStringPositionalArgument(&positionalArg, "username", "Another username argument")
	if err == nil {
		t.Fatalf("Expected an error for duplicate positional argument name, got none")
	}
}

func TestRegisterMultiplePositionalArguments(t *testing.T) {
	ap := &ArgumentsParser{}
	ap.PositionalArguments = make([]positionals.PositionalArgument, 0)

	var arg1 string
	err := ap.NewStringPositionalArgument(&arg1, "arg1", "First positional argument")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var arg2 string
	err = ap.NewStringPositionalArgument(&arg2, "arg2", "Second positional argument")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(ap.PositionalArguments) != 2 {
		t.Errorf("Expected 2 positional arguments, got %d", len(ap.PositionalArguments))
	}

	posarg0 := ap.PositionalArguments[0]
	posarg1 := ap.PositionalArguments[1]

	if posarg0.GetName() != "arg1" {
		t.Errorf("Expected positional argument name 'arg1', got %s", posarg0.GetName())
	}

	if posarg1.GetName() != "arg2" {
		t.Errorf("Expected positional argument name 'arg2', got %s", posarg1.GetName())
	}
}
