# Integer Argument in goopts

The Integer argument type in the `goopts` library allows you to define command-line arguments that expect an integer value. This is useful for scenarios where numerical input is required, such as specifying counts, limits, or other numeric configurations.

## Overview

- **Argument Name**: You can specify an integer argument using a flag (e.g., `-n` or `--number`).
- **Default Value**: You can set a default value for the integer argument.
- **Usage**: The value provided by the user will be parsed as an integer and can be utilized within your application logic.

## Defining an Integer Argument

To define an Integer argument using `goopts`, you can do the following:

```go
var count int

ap.NewIntegerArgument(&count, "-n", "--number", 0, false, "Number of items to process.")
```

## How to Use

When running your program, you can include the integer argument followed by the desired value:

```
./your_program -n 5
```

In this example, the count variable will be set to 5, and you can use this value in your application logic.

## Example Usage

Here's a complete example demonstrating the use of an Integer argument:

```go
package main

import (
	"fmt"
	"github.com/p0dalirius/goopts/parser"
)

var count int

func main() {
	ap := parser.ArgumentsParser{Banner: "Example of Integer Argument using goopts"}

	// Define the Integer argument
	ap.NewIntegerArgument(&count, "-n", "--number", 0, false, "Number of items to process.")

	// Parse the arguments
	ap.Parse()

	fmt.Printf("Number of items to process: %d\n", count)
}
```

## Running the Example

To run the example, compile it and execute with the integer argument:

```
go build -o example_integer_argument
./example_integer_argument -n 10
```

When executed, it will output:

```
Number of items to process: 10
```

If you run it without specifying the integer argument, it will use the default value:

```
./example_integer_argument
```

The output will be:

```
Number of items to process: 0
```

This simple implementation demonstrates how to effectively use Integer arguments in your command-line applications with the goopts library.