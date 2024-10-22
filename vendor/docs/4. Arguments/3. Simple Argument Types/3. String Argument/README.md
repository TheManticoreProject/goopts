# String Argument in goopts

The String argument type in the `goopts` library allows you to define command-line arguments that expect a string value. This is useful for scenarios where textual input is required, such as specifying names, file paths, or other string configurations.

## Overview

- **Argument Name**: You can specify a string argument using a flag (e.g., `-s` or `--string`).
- **Default Value**: You can set a default value for the string argument.
- **Usage**: The value provided by the user will be parsed as a string and can be utilized within your application logic.

## Defining a String Argument

To define a String argument using `goopts`, you can do the following:

```go
var name string

ap.NewStringArgument(&name, "-s", "--string", "default", false, "Specify the name or string value.")
```

## How to Use

When running your program, you can include the string argument followed by the desired value:

```
./your_program -s "Hello, World!"
```

In this example, the name variable will be set to "Hello, World!", and you can use this value in your application logic.

## Example Usage

Here's a complete example demonstrating the use of a String argument:

```go
package main

import (
	"fmt"
	"github.com/p0dalirius/goopts/parser"
)

var name string

func main() {
	ap := parser.ArgumentsParser{Banner: "Example of String Argument using goopts"}

	// Define the String argument
	ap.NewStringArgument(&name, "-s", "--string", "default", false, "Specify the name or string value.")

	// Parse the arguments
	ap.Parse()

	fmt.Printf("Provided name: %s\n", name)
}
```

## Running the Example

To run the example, compile it and execute with the string argument:

```
go build -o example_string_argument
./example_string_argument -s "John Doe"
```

When executed, it will output:

```
Provided name: John Doe
```

If you run it without specifying the string argument, it will use the default value:

```
./example_string_argument
```

The output will be:

```
Provided name: default
```

This simple implementation demonstrates how to effectively use String arguments in your command-line applications with the goopts library.