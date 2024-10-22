# StringPositionalArgument Documentation

The `StringPositionalArgument` struct is part of the `goopts` library, designed to handle command-line positional arguments that represent string values. This documentation outlines its properties, methods, and how to use it effectively.

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

The `StringPositionalArgument` struct is utilized to capture string arguments from the command line, providing detailed help messages to guide users on how to use the argument.

## Struct Definition

```go
type StringPositionalArgument struct {
    Name     string
    Help     string
    Value    *string
    Required bool
}
```

## Properties

- `Name`: The name of the argument.
- `Help`: The help message that describes the purpose of the argument.
- `Value`: A pointer to a string variable where the parsed value will be stored.
- `Required`: Indicates whether this argument must be provided by the user.

## Methods

### GetName

```go
func (arg StringPositionalArgument) GetName() string
```

Returns the name of the argument.

### GetHelp

```go
func (arg StringPositionalArgument) GetHelp() string
```

Returns the help message for the argument.

### GetValue

```go
func (arg StringPositionalArgument) GetValue() any
```
Returns the current value of the argument as an any type.

### IsRequired

```go
func (arg StringPositionalArgument) IsRequired() bool
```
Returns true if the argument is required; otherwise, returns false.

### Init

```go
func (arg *StringPositionalArgument) Init(value *string, name string, help string)
```
Initializes the StringPositionalArgument with the specified parameters.

Parameters:
- `value`: A pointer to a string where the argument's value will be stored.
- `name`: The name of the argument.
- `help`: The help message describing the argument.

### Consume

```go
func (arg StringPositionalArgument) Consume(arguments []string) ([]string, error)
```

Processes the command-line arguments to set the value of the StringPositionalArgument.

Parameters:
- `arguments`: A slice of strings representing the command-line arguments.

Returns:
- A slice of strings representing the remaining arguments after processing the StringPositionalArgument.

