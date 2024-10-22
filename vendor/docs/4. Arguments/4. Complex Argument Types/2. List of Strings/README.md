# ListOfStringsArgument Documentation

The `ListOfStringsArgument` struct is part of the `goopts` library and provides functionality for handling arguments that accept multiple string values in command-line applications. This documentation describes its properties, methods, and usage.

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

## Overview

The `ListOfStringsArgument` is designed to handle command-line arguments that can accept multiple string values. It supports both short and long argument names and provides help messages for user guidance.

## Struct Definition

```go
type ListOfStringsArgument struct {
    ShortName     string
    LongName      string
    Help          string
    Value         *[]string
    DefaultValue  []string
    Required      bool
}
```

## Properties

- `ShortName`: The short version of the argument (e.g., `-l`).
- `LongName`: The long version of the argument (e.g., `--list`).
- `Help`: The help message describing the argument's purpose.
- `Value`: A pointer to a slice of strings where the parsed values will be stored.
- `DefaultValue`: The default value if the argument is not provided.
- `Required`: Indicates whether the argument is mandatory.

## Methods

GetShortName

```go
func (arg ListOfStringsArgument) GetShortName() string
```

Returns the short name of the argument.

### GetLongName

```go
func (arg ListOfStringsArgument) GetLongName() string
```

Returns the long name of the argument.

### GetHelp

```go
func (arg ListOfStringsArgument) GetHelp() string
```

Returns the help message for the argument.

### GetValue

```go
func (arg ListOfStringsArgument) GetValue() any
```

Returns the current values of the argument as an any type.

### GetDefaultValue

```go
func (arg ListOfStringsArgument) GetDefaultValue() any
```

Returns the default values of the argument as an any type.

### IsRequired

```go
func (arg ListOfStringsArgument) IsRequired() bool
```

Returns true if the argument is required; otherwise, false.

### Init

```go
func (arg *ListOfStringsArgument) Init(value *[]string, shortName, longName string, defaultValue []string, required bool, help string)
```

Initializes the ListOfStringsArgument with provided parameters.

Parameters:
- `value`: A pointer to a slice of strings where the parsed values will be stored.
- `shortName`: The short name of the argument (single character). If empty, it will be set to an empty string.
- `longName`: The long name of the argument (string). If empty, it will be set to an empty string.
- `defaultValue`: The default value of the argument.
- `help`: The help message describing the argument.

### Consume

```go
func (arg ListOfStringsArgument) Consume(arguments []string) ([]string, error)
```

Processes the command-line arguments and sets the values of the `ListOfStringsArgument`.

Parameters:
- `arguments`: A slice of strings representing the command-line arguments.

Returns:
- A slice of strings representing the remaining arguments after processing the `ListOfStringsArgument`.
- An error if the argument is required and not found.
