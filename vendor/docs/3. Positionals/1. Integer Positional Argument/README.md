# IntPositionalArgument Documentation

The `IntPositionalArgument` struct is part of the `goopts` library, designed to handle command-line positional arguments that represent integer values. This documentation outlines its properties, methods, and how to use it effectively.

## Table of Contents

- [Overview](#overview)
- [Struct Definition](#struct-definition)
- [Methods](#methods)
  - [GetName](#getname)
  - [GetHelp](#gethelp)
  - [GetValue](#getvalue)
  - [IsRequired](#isrequired)
  - [Init](#init)
  - [Consume](#consume)
- [Usage Example](#usage-example)

## Overview

The `IntPositionalArgument` struct is utilized to capture integer arguments from the command line, with support for both short and long flags. It also provides detailed help messages for better user guidance.

## Struct Definition

```go
type IntPositionalArgument struct {
    Name     string
    Help     string
    Value    *int
    Required bool
}
```

## Properties

- `Name`: The name of the argument.
- `Help`: The help message that describes the purpose of the argument.
- `Value`: A pointer to a string variable where the parsed value will be stored.
- `Required`: Indicates whether this argument must be provided by the user.

## Methods

### GetShortName

### GetName

```go
func (arg IntPositionalArgument) GetName() string
```

Returns the name of the argument.

### GetHelp

```go
func (arg IntPositionalArgument) GetHelp() string
```

Returns the help message for the argument.

### GetValue

```go
func (arg IntPositionalArgument) GetValue() any
```

Returns the current value of the argument as an any type.

### IsRequired

```go
func (arg IntPositionalArgument) IsRequired() bool
```

Returns true if the argument is required; otherwise, returns false.

### Init

```go
func (arg *IntPositionalArgument) Init(value *int64, shortName, longName string, defaultValue int64, help string)
```

Initializes the IntPositionalArgument with the specified parameters.

Parameters:
- `value`: A pointer to an int64 where the argument's value will be stored.
- `shortName`: The short name for the argument (single character). If empty, it will be set to an empty string.
- `longName`: The long name for the argument (string). If empty, it will be set to an empty string.
- `defaultValue`: The default value for the argument.
- `help`: The help message describing the argument.

### Consume

```go
func (arg IntPositionalArgument) Consume(arguments []string) ([]string, error)
```

Processes the command-line arguments to set the value of the IntPositionalArgument.

Parameters:
- `arguments`: A slice of strings representing the command-line arguments.

Returns:
- A slice of strings representing the remaining arguments after processing the IntPositionalArgument.
- An error if the argument is required and not found or if parsing fails.