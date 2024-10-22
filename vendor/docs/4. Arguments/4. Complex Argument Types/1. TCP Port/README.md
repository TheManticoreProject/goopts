# TcpPortArgument Documentation

The `TcpPortArgument` struct is part of the `goopts` library and provides functionality for handling TCP port arguments in command-line applications. This documentation describes its properties, methods, and usage.

## Table of Contents

- [Overview](#overview)
- [Struct Definition](#struct-definition)
- [Methods](#methods)
  - [GetShortName](#getshortname)
  - [GetLongName](#getlongname)
  - [GetHelp](#gethelp)
  - [GetValue](#getvalue)
  - [GetDefaultValue](#getdefaultvalue)
  - [IsRequired](#isrequired)
  - [Init](#init)
  - [Consume](#consume)
- [Usage Example](#usage-example)
- [Error Handling](#error-handling)

## Overview

The `TcpPortArgument` is designed to handle TCP port arguments, ensuring that values provided are within the valid range (0-65535). It supports both short and long argument names and provides help messages for user guidance.

## Struct Definition

```go
type TcpPortArgument struct {
    ShortName     string
    LongName      string
    Help          string
    Value         *int
    DefaultValue  int
    Required      bool
}
```

## Properties

- `ShortName`: The short version of the argument (e.g., `-p`).
- `LongName`: The long version of the argument (e.g., `--port`).
- `Help`: The help message describing the argument's purpose.
- `Value`: A pointer to the integer where the parsed value will be stored.
- `DefaultValue`: The default value if the argument is not provided.
- `Required`: Indicates whether the argument is mandatory.

## Methods

### GetShortName

```go
func (arg TcpPortArgument) GetShortName() string
```

Returns the short name of the argument.

### GetLongName

```go
func (arg TcpPortArgument) GetLongName() string
```

Returns the long name of the argument.

### GetHelp

```go
func (arg TcpPortArgument) GetHelp() string
```

Returns the help message for the argument, including the default value if it is not required.

### GetValue

```go
func (arg TcpPortArgument) GetValue() any
```

Returns the current value of the argument as an any type.

### GetDefaultValue

```go
func (arg TcpPortArgument) GetDefaultValue() any
```

Returns the default value of the argument as an any type.

### IsRequired

```go
func (arg TcpPortArgument) IsRequired() bool
```

Returns true if the argument is required; otherwise, false.

### Init

```go
func (arg *TcpPortArgument) Init(value *int, shortName, longName string, defaultValue int, required bool, help string)
```

Initializes the TcpPortArgument with provided parameters.

Parameters:
- `value`: A pointer to the integer where the parsed value will be stored.
- `shortName`: The short name of the argument.
- `longName`: The long name of the argument.
- `defaultValue`: The default value for the argument.
- `required`: A boolean indicating if the argument is required.
- `help`: The help message for the argument.

### Consume

```go
func (arg TcpPortArgument) Consume(arguments []string) ([]string, error)
```

Processes the command-line arguments and sets the value of the TcpPortArgument.

Parameters:
- `arguments`: A slice of strings representing the command-line arguments.

Returns:
- A slice of strings representing the remaining arguments after processing the TcpPortArgument.
- An error if parsing fails or if the port value is out of the valid range.