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
