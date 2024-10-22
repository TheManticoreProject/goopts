package utils

import (
	"testing"
)

func TestStripLeftDashes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"---example", "example"},
		{"--example", "example"},
		{"-example", "example"},
		{"example", "example"},
		{"", ""},
	}

	for _, test := range tests {
		result := StripLeftDashes(test.input)
		if result != test.expected {
			t.Errorf("StripLeftDashes(%q) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestListOfStrings(t *testing.T) {
	tests := []struct {
		input    []string
		expected string
	}{
		{[]string{"one", "two", "three"}, "[\"one\", \"two\", \"three\"]"},
		{[]string{"single"}, "[\"single\"]"},
		{[]string{}, "[]"},
	}

	for _, test := range tests {
		result := ListOfStrings(test.input)
		if result != test.expected {
			t.Errorf("ListOfStrings(%v) = %q; expected %q", test.input, result, test.expected)
		}
	}
}

func TestStringToInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		hasError bool
	}{
		// Decimal tests
		{"123456789", 123456789, false},
		{"0", 0, false},
		{"-42", -42, false},

		// Hexadecimal tests
		{"0xabcdef", 11259375, false},
		{"0XABCDEF", 11259375, false},
		{"0x0", 0, false},

		// Octal tests
		{"0o1234567", 342391, false},
		{"0O765", 501, false},
		{"0o0", 0, false},

		// Binary tests
		{"0b101010", 42, false},
		{"0B1101", 13, false},
		{"0b0", 0, false},

		// Invalid input tests
		{"abcdef", 0, true},
		{"0xGHIJK", 0, true},
		{"0o1238", 0, true}, // Invalid octal digit
		{"0b102", 0, true},  // Invalid binary digit
		{"", 0, true},
	}

	for _, test := range tests {
		result, err := StringToInt(test.input)
		if test.hasError {
			if err == nil {
				t.Errorf("Expected an error for input '%s', but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Did not expect an error for input '%s', but got: %v", test.input, err)
			}
			if result != test.expected {
				t.Errorf("For input '%s', expected %d but got %d", test.input, test.expected, result)
			}
		}
	}
}
