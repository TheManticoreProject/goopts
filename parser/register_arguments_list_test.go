package parser

import (
	"testing"

	"github.com/TheManticoreProject/goopts/argumentgroup"
)

func TestRegisterListOfIntsArgument(t *testing.T) {
	ap := &ArgumentsParser{}
	ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)

	var listArg []int
	err := ap.NewListOfIntsArgument(&listArg, "-l", "--list", []int{1, 2, 3}, false, "List of integers argument")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(ap.Groups[""].Arguments) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(ap.Groups[""].Arguments))
	}

	if _, exists := ap.Groups[""].ShortNameToArgument["-l"]; !exists {
		t.Errorf("Expected short name '-l' to be registered")
	}

	if _, exists := ap.Groups[""].LongNameToArgument["--list"]; !exists {
		t.Errorf("Expected long name '--list' to be registered")
	}

	// Test duplicate registration
	err = ap.NewListOfIntsArgument(&listArg, "-l", "--list", []int{4, 5, 6}, false, "Another list of integers argument")
	if err == nil {
		t.Fatalf("Expected an error for duplicate short name, got none")
	}
}
