# Basic Example of Setting Up goopts

In this section, we will provide a simple example of how to set up and use the `goopts` library in your Go application.

## Step 1: Import the Library

First, ensure that you have the `goopts` library installed. Then, import it into your Go program:

```go
import (
    "fmt"
    "github.com/p0dalirius/goopts/parser"
)
```

## Step 2: Define Your Arguments

Next, you need to define the arguments that your application will accept. For this example, we will create a simple application that takes a URL and an optional verbosity flag.

```go
var (
    // Positional argument
    url string

    // Verbosity flag
    verbose bool
)
```

## Step 3: Create the Argument Parser

Create an instance of the `ArgumentsParser` and define your arguments:

```go
func parseArgs() {
    ap := parser.ArgumentsParser{Banner: "Simple goopts Example - by @podalirius"}

    // Define a positional argument
    ap.NewStringPositionalArgument(&url, "url", "URL to process.")

    // Define a boolean flag
    ap.NewBoolArgument(&verbose, "-v", "--verbose", false, "Enable verbose output.")

    // Parse the arguments
    ap.Parse()
}
```

## Step 4: Use the Parsed Arguments

In the `main` function, call `parseArgs()` and utilize the parsed arguments:

```go
func main() {
    parseArgs()

    // Use the parsed URL and verbosity flag
    fmt.Printf("Processing URL: %s\n", url)
    if verbose {
        fmt.Println("Verbose mode enabled.")
    }
}
```

## Step 5: Running Your Application

To run your application, compile it and pass the required arguments:

```bash
go run main.go http://example.com -v
```

You should see the following output:

```
Processing URL: http://example.com
Verbose mode enabled.
```

## Conclusion

You have now set up a basic application using the `goopts` library! You can extend this example by adding more arguments and functionality as needed. For further details, check out the [Documentation](https://github.com/p0dalirius/goopts/docs).
