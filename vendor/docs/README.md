![](../../.github/banner.png)

# goopts: A Command-Line Argument Parsing Library for Go

## Introduction

goopts is a powerful and flexible library designed to simplify the parsing of command-line arguments in Go applications. Inspired by popular argument parsing libraries from other languages, goopts aims to provide developers with a robust, easy-to-use interface for handling various argument types, including positional arguments, named options, flags, and more. The library also supports advanced features like argument groups, mutually exclusive options, and dependency handling, making it an excellent choice for building feature-rich CLI tools.

Whether you are developing a simple utility or a complex command-line application, goopts offers the tools you need to parse and manage user input efficiently and effectively.

## Purpose

The purpose of goopts is to:

- Provide a straightforward and intuitive way to define and parse command-line arguments.
- Support multiple types of arguments including strings, integers, booleans, lists, and maps.
- Enable the grouping of arguments, including mutually exclusive groups and dependent arguments.
- Offer clear, customizable help messages for your command-line tools.
- Allow for flexible configuration and easy integration into existing Go projects.

## Table of Contents

- [1. Getting Started](./1.%20Getting%20Started/README.md)
    - [1.1 Installation Instructions](./1.%20Getting%20Started/1.%20Installation%20Instructions/README.md)
    - [1.2 Basic Example of Setting Up goopts](./1.%20Getting%20Started/2.%20Basic%20Example%20of%20Setting%20Up%20goopts/README.md)
    - [1.3 Positional vs. Named Arguments](./1.%20Getting%20Started/3.%20Positional%20vs.%20Named%20Arguments/README.md)
- [2. Examples](./2.%20Examples/README.md)
    - [2.1. Basic Example](./2.%20Examples/1.%20Basic%20Example/README.md)
    - [2.2. Parsing HTTP Headers and Other Advanced Use Cases](./2.%20Examples/2.%20Parsing%20HTTP%20Headers%20and%20Other%20Advanced%20Use%20Cases/README.md)
    - [2.3. Complex CLI Tool with Multiple Groups](./2.%20Examples/3.%20Complex%20CLI%20Tool%20with%20Multiple%20Groups/README.md)
- [3. Positionals](./3.%20Positionals/README.md)
    - [3.1. Integer Argument](./3.%20Positionals/1.%20Integer%20Positional%20Argument/README.md)
    - [3.2. String Argument](./3.%20Positionals/2.%20String%20Positional%20Argument/README.md)
- [4. Arguments](./4.%20Arguments/README.md)
    - [3. Simple Argument Types](./4.%20Arguments/3.%20Simple%20Argument%20Types/README.md)
        - [4.3.1. Boolean Argument](./4.%20Arguments/3.%20Simple%20Argument%20Types/1.%20Boolean%20Argument/README.md)
        - [4.3.2. Integer Argument](./4.%20Arguments/3.%20Simple%20Argument%20Types/2.%20Integer%20Argument/README.md)
        - [4.3.3. String Argument](./4.%20Arguments/3.%20Simple%20Argument%20Types/3.%20String%20Argument/README.md)
    - [4. Complex Argument Types](./4.%20Arguments/4.%20Complex%20Argument%20Types/README.md)
        - [4.4.1. TCP Port](./4.%20Arguments/4.%20Complex%20Argument%20Types/1.%20TCP%20Port/README.md)
        - [4.4.2. List of Strings](./4.%20Arguments/4.%20Complex%20Argument%20Types/2.%20List%20of%20Strings/README.md)
        - [4.4.3. Map of HTTP Headers](./4.%20Arguments/4.%20Complex%20Argument%20Types/3.%20Map%20of%20HTTP%20Headers/README.md)
- [5. Argument Groups](./5.%20Argument%20Groups/README.md)
    - [1. Mutually Exclusive Argument Groups](./5.%20Argument%20Groups/1.%20Mutually%20Exclusive%20Argument%20Groups/README.md)
    - [2. Dependent Argument Groups](./5.%20Argument%20Groups/2.%20Dependent%20Argument%20Groups/README.md)
- [6. Best Practices](./6.%20Best%20Practices/README.md)
    - [1. Organizing Your Code with goopts](./6.%20Best%20Practices/1.%20Organizing%20Your%20Code%20with%20goopts/README.md)
    - [2. Common Pitfalls and How to Avoid Them](./6.%20Best%20Practices/2.%20Common%20Pitfalls%20and%20How%20to%20Avoid%20Them/README.md)
