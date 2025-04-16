package parser

import (
	"testing"

	"github.com/TheManticoreProject/goopts/argumentgroup"
)

func TestRegisterBoolArgument(t *testing.T) {
	ap := &ArgumentsParser{}
	ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)

	var boolArg bool
	err := ap.NewBoolArgument(&boolArg, "-b", "--bool", false, "A boolean argument")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(ap.Groups[""].Arguments) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(ap.Groups[""].Arguments))
	}

	if _, exists := ap.Groups[""].ShortNameToArgument["-b"]; !exists {
		t.Errorf("Expected short name '-b' to be registered")
	}

	if _, exists := ap.Groups[""].LongNameToArgument["--bool"]; !exists {
		t.Errorf("Expected long name '--bool' to be registered")
	}

	// Test duplicate registration
	err = ap.NewBoolArgument(&boolArg, "-b", "--bool", true, "Another boolean argument")
	if err == nil {
		t.Fatalf("Expected an error for duplicate short name, got none")
	}
}

func TestRegisterStringArgument(t *testing.T) {
	ap := &ArgumentsParser{}
	ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)

	var stringArg string
	err := ap.NewStringArgument(&stringArg, "-u", "--username", "defaultUser", false, "Username argument")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(ap.Groups[""].Arguments) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(ap.Groups[""].Arguments))
	}

	if _, exists := ap.Groups[""].ShortNameToArgument["-u"]; !exists {
		t.Errorf("Expected short name '-u' to be registered")
	}

	if _, exists := ap.Groups[""].LongNameToArgument["--username"]; !exists {
		t.Errorf("Expected long name '--username' to be registered")
	}

	// Test duplicate registration
	err = ap.NewStringArgument(&stringArg, "-u", "--username", "newUser", false, "Another username argument")
	if err == nil {
		t.Fatalf("Expected an error for duplicate short name, got none")
	}
}

func TestRegisterIntArgument(t *testing.T) {
	ap := &ArgumentsParser{}
	ap.Groups = make(map[string]*argumentgroup.ArgumentGroup)

	var intArg int
	err := ap.NewIntArgument(&intArg, "-i", "--integer", 42, false, "Integer argument")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(ap.Groups[""].Arguments) != 1 {
		t.Errorf("Expected 1 argument, got %d", len(ap.Groups[""].Arguments))
	}

	if _, exists := ap.Groups[""].ShortNameToArgument["-i"]; !exists {
		t.Errorf("Expected short name '-i' to be registered")
	}

	if _, exists := ap.Groups[""].LongNameToArgument["--integer"]; !exists {
		t.Errorf("Expected long name '--integer' to be registered")
	}

	// Test duplicate registration
	err = ap.NewIntArgument(&intArg, "-i", "--integer", 99, false, "Another integer argument")
	if err == nil {
		t.Fatalf("Expected an error for duplicate short name, got none")
	}
}
