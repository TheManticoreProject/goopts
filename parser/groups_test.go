package parser

import (
	"testing"
)

func TestNewArgumentGroup(t *testing.T) {
	ap := &ArgumentsParser{}

	groupName := "testGroup"
	group, err := ap.NewArgumentGroup(groupName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if group.Name != groupName {
		t.Errorf("Expected group name '%s', got '%s'", groupName, group.Name)
	}

	if ap.Groups[groupName] != group {
		t.Errorf("Expected group to be stored in ap.Groups, but it wasn't")
	}

	// Test duplicate group creation
	_, err = ap.NewArgumentGroup(groupName)
	if err == nil {
		t.Fatalf("Expected an error for duplicate group name, got none")
	}
}

func TestNewNotRequiredMutuallyExclusiveArgumentGroup(t *testing.T) {
	ap := &ArgumentsParser{}

	groupName := "notRequiredGroup"
	group, err := ap.NewNotRequiredMutuallyExclusiveArgumentGroup(groupName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if group.Name != groupName {
		t.Errorf("Expected group name '%s', got '%s'", groupName, group.Name)
	}

	if ap.Groups[groupName] != group {
		t.Errorf("Expected group to be stored in ap.Groups, but it wasn't")
	}

	// Test duplicate group creation
	_, err = ap.NewNotRequiredMutuallyExclusiveArgumentGroup(groupName)
	if err == nil {
		t.Fatalf("Expected an error for duplicate group name, got none")
	}
}

func TestNewRequiredMutuallyExclusiveArgumentGroup(t *testing.T) {
	ap := &ArgumentsParser{}

	groupName := "requiredGroup"
	group, err := ap.NewRequiredMutuallyExclusiveArgumentGroup(groupName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if group.Name != groupName {
		t.Errorf("Expected group name '%s', got '%s'", groupName, group.Name)
	}

	if ap.Groups[groupName] != group {
		t.Errorf("Expected group to be stored in ap.Groups, but it wasn't")
	}

	// Test duplicate group creation
	_, err = ap.NewRequiredMutuallyExclusiveArgumentGroup(groupName)
	if err == nil {
		t.Fatalf("Expected an error for duplicate group name, got none")
	}
}
