# Installation Instructions

To start using the `goopts` library, follow the steps below to install it and set it up in your Go project.

## Prerequisites

Make sure you have the following installed on your system:
- Go 1.18 or later
- Git (to clone the repository if installing from source)

## Installing via `go get`

The easiest way to install `goopts` is by using the `go get` command. Run the following command in your terminal:

```bash
go get github.com/p0dalirius/goopts
```

This will download and install the `goopts` library, making it available for import in your Go projects.

## Importing `goopts` in Your Project

After installation, you can import `goopts` in your Go project as follows:

```go
import (
    "github.com/p0dalirius/goopts/parser"
)
```

## Installing from Source

If you prefer to clone the repository and install the library from source, follow these steps:

1. Clone the repository:

    ```bash
    git clone https://github.com/p0dalirius/goopts.git
    ```

2. Navigate to the cloned directory:

    ```bash
    cd goopts
    ```

3. Install the library:

    ```bash
    go install
    ```

This will install `goopts` in your `$GOPATH/bin`, and you can use it like any other Go package.

## Verifying Installation

To verify that the installation was successful, you can run a basic `go build` or `go run` command with a simple Go program that imports `goopts`. For example:

```go
package main

import (
    "fmt"
    "github.com/p0dalirius/goopts/parser"
)

func main() {
    fmt.Println("Successfully installed goopts!")
}
```

Run the program:

```bash
go run main.go
```

If everything is set up correctly, you should see:

```
Successfully installed goopts!
```

## Updating `goopts`

To update `goopts` to the latest version, run:

```bash
go get -u github.com/p0dalirius/goopts
```

This command will fetch the latest version of the library and update your project.
