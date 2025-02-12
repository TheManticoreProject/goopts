package arguments

import (
	"reflect"
	"testing"
)

func TestMapOfHttpHeadersArgument_Init(t *testing.T) {
	var value map[string]string
	arg := MapOfHttpHeadersArgument{}
	arg.Init(&value, "h", "headers", map[string]string{"Content-Type": "application/json"}, true, "Help message for headers")
	arg.ResetDefaultValue()

	if arg.ShortName != "-h" {
		t.Errorf("Expected short name '-h', got '%s'", arg.ShortName)
	}

	if arg.LongName != "--headers" {
		t.Errorf("Expected long name '--headers', got '%s'", arg.LongName)
	}

	if arg.Help != "Help message for headers" {
		t.Errorf("Expected help message 'Help message for headers', got '%s'", arg.Help)
	}

	if !arg.Required {
		t.Errorf("Expected required true, got false")
	}

	if !reflect.DeepEqual(arg.DefaultValue, map[string]string{"Content-Type": "application/json"}) {
		t.Errorf("Expected default value '%v', got '%v'", map[string]string{"Content-Type": "application/json"}, arg.DefaultValue)
	}
}

func TestMapOfHttpHeadersArgument_Consume(t *testing.T) {
	value := make(map[string]string)
	arg := MapOfHttpHeadersArgument{
		ShortName:    "-h",
		LongName:     "--headers",
		Value:        &value,
		DefaultValue: map[string]string{},
		Required:     true,
	}
	arg.ResetDefaultValue()

	args := []string{"-h", "Authorization: Bearer token"}

	remainingArgs, err := arg.Consume(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(remainingArgs) != 0 {
		t.Errorf("Expected no remaining arguments, got '%v'", remainingArgs)
	}

	if value["Authorization"] != "Bearer token" {
		t.Errorf("Expected header 'Authorization: Bearer token', got '%s: %s'", "Authorization", value["Authorization"])
	}
}

func TestMapOfHttpHeadersArgument_Consume_MissingValue(t *testing.T) {
	value := make(map[string]string)
	arg := MapOfHttpHeadersArgument{
		ShortName:    "-h",
		LongName:     "--headers",
		Value:        &value,
		DefaultValue: map[string]string{},
		Required:     true,
	}
	arg.ResetDefaultValue()

	args := []string{"-h", "Authorization"}

	remainingArgs, err := arg.Consume(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(remainingArgs) != 0 {
		t.Errorf("Expected no remaining arguments, got '%v'", remainingArgs)
	}

	if value["Authorization"] != "" {
		t.Errorf("Expected empty value for 'Authorization', got '%s'", value["Authorization"])
	}
}

func TestMapOfHttpHeadersArgument_Consume_NoMatch(t *testing.T) {
	value := make(map[string]string)
	arg := MapOfHttpHeadersArgument{
		ShortName:    "-h",
		LongName:     "--headers",
		Value:        &value,
		DefaultValue: map[string]string{},
		Required:     true,
	}
	arg.ResetDefaultValue()

	args := []string{"--other"}

	remainingArgs, err := arg.Consume(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(remainingArgs) != 1 || remainingArgs[0] != "--other" {
		t.Errorf("Expected remaining arguments '--other', got '%v'", remainingArgs)
	}

	if len(value) != 0 {
		t.Errorf("Expected no headers to be set, got '%v'", value)
	}
}
