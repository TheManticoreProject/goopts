# CLI Tool with Multiple Modes and Subparsers Using goopts

This example demonstrates how to create a command-line application with multiple modes and subparsers using the `goopts` library. It showcases the ability to define different modes of operation, each with its own set of arguments, allowing for a versatile and organized command-line interface.

## Program Description

### Arguments

- **Positional Arguments**:
  - `filePath (string)`: The file to process.
  - `outputFolder (string)`: The output folder for results.

- **Recon Mode Arguments**:
  - `-f`, `--file (string)`: The file to recon.
  - `-o`, `--output (string)`: The output folder.
  - `-t`, `--timeout (int)`: Timeout for the recon operation.
  - `-v`, `--verbose (boolean)`: Enable verbose output for recon.

- **Scan Mode Arguments**:
  - `-f`, `--file (string)`: The file to scan.
  - `-d`, `--depth (int)`: Depth of the scan.
  - `-q`, `--quick (boolean)`: Enable quick scan mode.

### How to Run

Compile and run the example program using the following commands:

```bash
go build -o goopts_example main.go
./goopts_example recon -f file.txt -o output -t 60 -v
./goopts_example scan -f file.txt -d 5 -q
```

This will demonstrate the program's ability to handle different modes and subparsers, each with its own set of arguments.