# Basic Example Using goopts Parsing Library

This example demonstrates a simple command-line application using the `goopts` library to parse basic arguments. It will show how to set up an `ArgumentsParser`, define a few common argument types, and handle user input.

## Program Description

### Positional Arguments

- `name (string)`: A required positional argument to specify a user's name.

### Flags

- `-a`, `--age (int)`: Optional flag to specify the user's age.
- `-v`, `--verbose (boolean)`: Optional flag to enable verbose mode. Default is `false`.

### How to Run

Compile and run the example program using the following commands:

```bash
go build -o basic_example
./basic_example <name> [options]
```

Example Execution

```bash
./basic_example Alice -a 30 -v
```

When executed, the program will display:

```csharp
Hello, Alice!
Your age is: 30
Verbose mode is enabled.
```

If you run without the `-v` flag:

```bash
./basic_example Alice -a 30
```

It will display:

```
Hello, Alice!
Your age is: 30
```

## Code Walkthrough

The program starts by setting up an ArgumentsParser:

```go
package main

import (
    "fmt"
    "log"
    
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
    ap.NewStringPositionalArgument(&name, "name", "Your name.")

    // Define optional flags
    ap.NewIntArgument(&age, "-a", "--age", 0, false, "Specify your age.")
    ap.NewBoolArgument(&verbose, "-v", "--verbose", false, "Enable verbose mode.")

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
```

## Explanation

1. Positional Argument:
    ```go
    ap.NewStringPositionalArgument(&name, "name", "Your name.")
    ```
    This defines a required positional argument name that must be provided by the user.

2. Named Arguments:
    - `-a`, `--age`: Optional flag to specify the user's age.
        ```go
        ap.NewIntArgument(&age, "-a", "--age", 0, false, "Specify your age.")
        ```

    - `-v`, `--verbose`: Optional flag to enable verbose mode.
        ```go
        ap.NewBoolArgument(&verbose, "-v", "--verbose", false, "Enable verbose mode.")
        ```

With this simple setup, goopts allows you to quickly parse user inputs and handle optional and positional arguments effortlessly. This is a basic introduction to the features of goopts, showing how easy it is to build command-line tools.