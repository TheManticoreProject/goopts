package main

import (
	"fmt"

	"github.com/p0dalirius/goopts/parser"
)

var (
	// Positional argument
	name string

	// Optional flags
	age     int
	verbose bool
)

func parseArgs() {
	// Create a new arguments parser
	ap := parser.ArgumentsParser{Banner: "Basic Example of goopts v1.0"}

	// Define positional argument
	ap.NewStringPositionalArgument(&name, "name", "Your name")

	// Define optional flags
	ap.NewIntArgument(&age, "-a", "--age", 0, false, "Specify your age")
	ap.NewBoolArgument(&verbose, "-v", "--verbose", false, "Enable verbose mode")

	// Parse the arguments
	ap.Parse()
}

func main() {
	parseArgs()

	// Print greeting
	fmt.Printf("Hello, %s!\n", name)
	if age > 0 {
		fmt.Printf("Your age is: %d\n", age)
	}

	// Check verbose flag
	if verbose {
		fmt.Println("Verbose mode is enabled.")
	}
}
