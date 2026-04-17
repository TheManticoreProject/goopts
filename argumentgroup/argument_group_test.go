package argumentgroup

import (
	"testing"
)

func TestArgumentGroup_Register_ExistingArgument(t *testing.T) {
	var intValue1, intValue2 int
	argGroup := ArgumentGroup{Name: "Test Group"}

	// First registration should succeed
	err := argGroup.NewIntArgument(&intValue1, "i", "iterations", 10, true, "Number of iterations")
	if err != nil {
		t.Fatalf("Expected no error on first registration, got %v", err)
	}

	// Second registration with the same short name should fail
	err = argGroup.NewIntArgument(&intValue2, "i", "iterations", 5, true, "Another description")
	if err == nil {
		t.Fatal("Expected error on second registration, got none")
	}
}

func TestArgumentGroup_Register_ShortNameOnlyDoesNotPolluteLongNameMap(t *testing.T) {
	var v1, v2 bool
	ag := ArgumentGroup{Name: "G"}

	if err := ag.NewBoolArgument(&v1, "x", "", false, "help"); err != nil {
		t.Fatalf("unexpected error registering first arg: %v", err)
	}
	if _, ok := ag.LongNameToArgument[""]; ok {
		t.Fatalf("LongNameToArgument should not contain an empty-string key, got %v", ag.LongNameToArgument)
	}

	if err := ag.NewBoolArgument(&v2, "x", "", false, "help"); err == nil {
		t.Fatalf("duplicate short name should be rejected, got nil error")
	}
}

func TestArgumentGroup_Register_TwoShortNameOnlyArgsCoexist(t *testing.T) {
	var v1, v2 bool
	ag := ArgumentGroup{Name: "G"}

	if err := ag.NewBoolArgument(&v1, "x", "", false, "help"); err != nil {
		t.Fatalf("unexpected error registering first arg: %v", err)
	}
	if err := ag.NewBoolArgument(&v2, "y", "", false, "help"); err != nil {
		t.Fatalf("second short-only arg must register without error, got %v", err)
	}
	if len(ag.Arguments) != 2 {
		t.Fatalf("expected 2 registered arguments, got %d", len(ag.Arguments))
	}
	if len(ag.LongNameToArgument) != 0 {
		t.Fatalf("LongNameToArgument should be empty, got %v", ag.LongNameToArgument)
	}
}

func TestArgumentGroup_PrintArgumentTree(t *testing.T) {
	var stringValue string
	argGroup := ArgumentGroup{Name: "Test Group"}
	argGroup.NewStringArgument(&stringValue, "f", "file", "default.txt", true, "File to process")

	// Capture the output
	// NOTE: This can be enhanced by using a custom writer to capture output for validation
	argGroup.PrintArgumentTree(0)

	// This test simply checks that the function runs without panic.
}

func TestArgumentGroup_PrintArgumentTree_Empty(t *testing.T) {
	argGroup := ArgumentGroup{Name: "Empty Group"}
	argGroup.PrintArgumentTree(0)

	// This test simply checks that the function runs without panic for an empty group.
}
