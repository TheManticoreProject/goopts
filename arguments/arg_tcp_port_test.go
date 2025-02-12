package arguments

import (
	"fmt"
	"testing"
)

func TestTcpPortArgument_Init(t *testing.T) {
	var value int
	arg := TcpPortArgument{}
	arg.Init(&value, "p", "port", 8080, true, "Port to connect to")
	arg.ResetDefaultValue()

	if arg.ShortName != "-p" {
		t.Errorf("Expected short name '-p', got '%s'", arg.ShortName)
	}

	if arg.LongName != "--port" {
		t.Errorf("Expected long name '--port', got '%s'", arg.LongName)
	}

	if arg.Help != "Port to connect to" {
		t.Errorf("Expected help message 'Port to connect to', got '%s'", arg.Help)
	}

	if !arg.Required {
		t.Errorf("Expected required to be true, got false")
	}

	if arg.DefaultValue != 8080 {
		t.Errorf("Expected default value '8080', got '%d'", arg.DefaultValue)
	}
}

func TestTcpPortArgument_Consume_ValidPort(t *testing.T) {
	var value int
	arg := TcpPortArgument{
		ShortName:    "-p",
		LongName:     "--port",
		Value:        &value,
		DefaultValue: 8080,
		Required:     true,
	}
	arg.ResetDefaultValue()

	args := []string{"-p", "9090"}

	remainingArgs, err := arg.Consume(args)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(remainingArgs) != 0 {
		t.Errorf("Expected no remaining arguments, got '%v'", remainingArgs)
	}

	if value != 9090 {
		t.Errorf("Expected value '9090', got '%d'", value)
	}
}

func TestTcpPortArgument_Consume_InvalidPort(t *testing.T) {
	var value int
	arg := TcpPortArgument{
		ShortName:    "-p",
		LongName:     "--port",
		Value:        &value,
		DefaultValue: 8080,
		Required:     true,
	}
	arg.ResetDefaultValue()

	args := []string{"-p", "70000"}

	_, err := arg.Consume(args)
	if err == nil {
		t.Error("Expected error for port out of range, got nil")
	}

	expectedError := fmt.Sprintf("%s %s: TCP port value has to be in range 0-65535", args[0], args[1])
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', got '%v'", expectedError, err)
	}
}

func TestTcpPortArgument_Consume_InvalidValue(t *testing.T) {
	var value int
	arg := TcpPortArgument{
		ShortName:    "-p",
		LongName:     "--port",
		Value:        &value,
		DefaultValue: 8080,
		Required:     true,
	}
	arg.ResetDefaultValue()

	args := []string{"-p", "invalid"}

	_, err := arg.Consume(args)
	if err == nil {
		t.Error("Expected error for invalid integer value, got nil")
	}

	expectedError := fmt.Sprintf("%s %s: could not parse integer", args[0], args[1])
	if err.Error()[:len(expectedError)] != expectedError {
		t.Errorf("Expected error message '%s', got '%v'", expectedError, err)
	}
}

func TestTcpPortArgument_Consume_NoMatch(t *testing.T) {
	var value int
	arg := TcpPortArgument{
		ShortName:    "-p",
		LongName:     "--port",
		Value:        &value,
		DefaultValue: 8080,
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

	if value != 0 {
		t.Errorf("Expected value to remain unset, got '%d'", value)
	}
}
