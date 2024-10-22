# Argument Groups in goopts

This section discusses the concept of argument groups in the `goopts` argument parsing library. Argument groups allow for the organization of related arguments, enabling more structured command-line applications.

## Table of Contents

- [1. Mutually Exclusive Argument Groups](./1.%20Mutually%20Exclusive%20Argument%20Groups/README.md)
- [2. Dependent Argument Groups](./2.%20Dependent%20Arguments/README.md)

## Overview of Argument Groups

- **Mutually Exclusive Argument Groups**: This group type ensures that only one argument from a set can be provided by the user. It is useful for situations where certain options are alternatives, preventing conflicting inputs.

- **Dependent Argument Groups**: This group type establishes dependencies between arguments, requiring certain arguments to be provided only if others are specified. This is helpful for creating complex command-line interfaces where the relationships between options matter.

By utilizing argument groups, you can create more coherent command-line applications that guide user input and enforce logical relationships between options.
