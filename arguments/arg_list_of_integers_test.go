package arguments

import (
	"fmt"
	"slices"
	"testing"
)

func TestListOfIntsArgument_Init(t *testing.T) {
	var values []int

	arg := ListOfIntsArgument{}
	arg.Init(&values, "n", "numbers", []int{1, 2, 3}, true, "List of numbers")
	arg.ResetDefaultValue()

	if arg.ShortName != "-n" {
		t.Errorf("Expected ShortName to be '-n', got '%s'", arg.ShortName)
	}
	if arg.LongName != "--numbers" {
		t.Errorf("Expected LongName to be '--numbers', got '%s'", arg.LongName)
	}
	if !arg.Required {
		t.Errorf("Expected Required to be true, got '%v'", arg.Required)
	}
	if arg.Help != "List of numbers" {
		t.Errorf("Expected Help to be 'List of numbers', got '%s'", arg.Help)
	}
	if len(arg.DefaultValue) != 3 || arg.DefaultValue[0] != 1 || arg.DefaultValue[1] != 2 || arg.DefaultValue[2] != 3 {
		t.Errorf("Expected DefaultValue to be [1, 2, 3], got '%v'", arg.DefaultValue)
	}
	if !slices.Equal(*arg.Value, arg.DefaultValue) {
		t.Errorf("Expected Value to be equal to DefaultValue, got '%v'", *arg.Value)
	}
}

func TestListOfIntsArgument_Consume(t *testing.T) {
	values := []int{}
	arg := ListOfIntsArgument{
		ShortName: "-n",
		LongName:  "--numbers",
		Value:     &values,
	}
	arg.ResetDefaultValue()

	arguments := []string{"-n", "42", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if len(*arg.Value) != 1 || (*arg.Value)[0] != 42 {
		t.Errorf("Expected Value to contain [42], got '%v'", *arg.Value)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "anotherArg" {
		t.Errorf("Expected remaining arguments to be '[anotherArg]', got '%v'", remainingArgs)
	}
}

func TestListOfIntsArgument_Consume_Multiple(t *testing.T) {
	values := []int{}
	arg := ListOfIntsArgument{
		ShortName: "-n",
		LongName:  "--numbers",
		Value:     &values,
	}
	arg.ResetDefaultValue()

	arguments := []string{"-n", "42", "-n", "100", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)
	remainingArgs, _ = arg.Consume(remainingArgs)

	if len(*arg.Value) != 2 || (*arg.Value)[0] != 42 || (*arg.Value)[1] != 100 {
		t.Errorf("Expected Value to contain [42, 100], got '%v'", *arg.Value)
	}
	if len(remainingArgs) != 1 || remainingArgs[0] != "anotherArg" {
		t.Errorf("Expected remaining arguments to be '[anotherArg]', got '%v'", remainingArgs)
	}
}

func TestListOfIntsArgument_Consume_NoMatch(t *testing.T) {
	values := []int{}
	arg := ListOfIntsArgument{
		ShortName: "-n",
		LongName:  "--numbers",
		Value:     &values,
	}
	arg.ResetDefaultValue()

	arguments := []string{"-x", "42", "anotherArg"}
	remainingArgs, _ := arg.Consume(arguments)

	if len(*arg.Value) != 0 {
		t.Errorf("Expected Value to remain empty, got '%v'", *arg.Value)
	}
	if len(remainingArgs) != 3 {
		t.Errorf("Expected remaining arguments to be the same as input, got '%v'", remainingArgs)
	}
}

func TestListOfIntsArgument_Getters(t *testing.T) {
	values := []int{7, 14}
	arg := ListOfIntsArgument{
		ShortName:    "-n",
		LongName:     "--numbers",
		Help:         "List of integers",
		Value:        &values,
		DefaultValue: []int{1, 2, 3},
		Required:     true,
	}
	arg.ResetDefaultValue()

	if arg.GetShortName() != "-n" {
		t.Errorf("Expected ShortName to be '-n', got '%s'", arg.GetShortName())
	}
	if arg.GetLongName() != "--numbers" {
		t.Errorf("Expected LongName to be '--numbers', got '%s'", arg.GetLongName())
	}
	if arg.GetHelp() != "List of integers" {
		t.Errorf("Expected Help to be 'List of integers', got '%s'", arg.GetHelp())
	}
	if val, ok := arg.GetValue().([]int); ok {
		if !slices.Equal(val, values) {
			fmt.Printf("Expected Value to be '%v', got '%v'\n", values, val)
		}
	} else {
		fmt.Printf("Expected GetValue to return a []int, but got a different type\n")
	}
	if val, ok := arg.GetDefaultValue().([]int); ok {
		if !slices.Equal(val, []int{1, 2, 3}) {
			fmt.Printf("Expected DefaultValue to be [1, 2, 3], got '%v'\n", val)
		}
	} else {
		fmt.Printf("Expected GetDefaultValue to return a []int, but got a different type\n")
	}
	if !arg.IsRequired() {
		t.Errorf("Expected Required to be true, got '%v'", arg.IsRequired())
	}
}
