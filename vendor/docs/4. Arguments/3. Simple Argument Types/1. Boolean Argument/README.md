# Boolean Argument in goopts

The Boolean argument type in the `goopts` library is used to represent a true/false value in command-line applications. It allows you to enable or disable features easily using flags. This can be particularly useful for toggling options in your CLI tool.

## Overview

- **Argument Name**: The name of the Boolean argument can be specified using a flag (e.g., `-v` or `--verbose`).
- **Default Value**: You can set a default value for the argument (true or false).
- **Usage**: When the flag is present in the command line, the argument is set to true; otherwise, it defaults to false.

## Defining a Boolean Argument

Here's how you can define a Boolean argument using `goopts`:

```go
var verbose bool

ap.NewBoolArgument(&verbose, "-v", "--verbose", false, "Enable verbose mode to see detailed output.")
```

## How to Use

When you run your program, you can include the Boolean flag to enable verbose output:

```
./your_program -v
```

In this example, the verbose variable will be set to true, and you can adjust the behavior of your application accordingly.

## Example Usage

Here's a complete example that uses a Boolean argument:

```go
package main

import (
	"fmt"
	"github.com/p0dalirius/goopts/parser"
)

var verbose bool

func main() {
	ap := parser.ArgumentsParser{Banner: "Example of Boolean Argument using goopts"}

	// Define the Boolean argument
	ap.NewBoolArgument(&verbose, "-v", "--verbose", false, "Enable verbose mode to see detailed output.")

	// Parse the arguments
	ap.Parse()

	if verbose {
		fmt.Println("Verbose mode enabled.")
	} else {
		fmt.Println("Verbose mode disabled.")
	}
}
```

## Running the Example

To run the example, compile it and execute with the -v flag:

```
go build -o example_boolean_argument
./example_boolean_argument -v
```

When executed, it will output:

```
Verbose mode enabled.
```

If you run it without the flag:

```
./example_boolean_argument
```

The output will be:

```
Verbose mode disabled.
```

This simple implementation demonstrates how to effectively use Boolean arguments in your command-line applications with the goopts library.