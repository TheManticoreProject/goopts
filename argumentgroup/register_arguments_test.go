package argumentgroup

import (
	"testing"
)

func TestArgumentGroup_NewIntArgument(t *testing.T) {
	var intValue int
	argGroup := ArgumentGroup{Name: "Test Group"}

	err := argGroup.NewIntArgument(&intValue, "i", "iterations", 10, true, "Number of iterations")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if argGroup.Arguments == nil || len(argGroup.Arguments) != 1 {
		t.Errorf("Expected 1 argument registered, got %d", len(argGroup.Arguments))
	}
	arg := argGroup.Arguments[0]
	if arg.GetShortName() != "-i" {
		t.Errorf("Expected short name to be '-i', got '%s'", arg.GetShortName())
	}
	if arg.GetLongName() != "--iterations" {
		t.Errorf("Expected long name to be '--iterations', got '%s'", arg.GetLongName())
	}
}

func TestArgumentGroup_NewBoolArgument(t *testing.T) {
	var boolValue bool
	argGroup := ArgumentGroup{Name: "Test Group"}

	err := argGroup.NewBoolArgument(&boolValue, "S", "use-ldaps", false, "Use LDAPS for secure connections")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if argGroup.Arguments == nil || len(argGroup.Arguments) != 1 {
		t.Errorf("Expected 1 argument registered, got %d", len(argGroup.Arguments))
	}
	arg := argGroup.Arguments[0]
	if arg.GetShortName() != "-S" {
		t.Errorf("Expected short name to be '-S', got '%s'", arg.GetShortName())
	}
	if arg.GetLongName() != "--use-ldaps" {
		t.Errorf("Expected long name to be '--use-ldaps', got '%s'", arg.GetLongName())
	}
}

func TestArgumentGroup_NewStringArgument(t *testing.T) {
	var stringValue string
	argGroup := ArgumentGroup{Name: "Test Group"}

	err := argGroup.NewStringArgument(&stringValue, "u", "username", "defaultUser", true, "Username for authentication")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if argGroup.Arguments == nil || len(argGroup.Arguments) != 1 {
		t.Errorf("Expected 1 argument registered, got %d", len(argGroup.Arguments))
	}
	arg := argGroup.Arguments[0]
	if arg.GetShortName() != "-u" {
		t.Errorf("Expected short name to be '-u', got '%s'", arg.GetShortName())
	}
	if arg.GetLongName() != "--username" {
		t.Errorf("Expected long name to be '--username', got '%s'", arg.GetLongName())
	}
}
