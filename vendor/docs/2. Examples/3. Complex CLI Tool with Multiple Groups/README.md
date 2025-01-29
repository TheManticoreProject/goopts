# Complex CLI Tool with Multiple Groups Using goopts

This example demonstrates how to create a more advanced command-line application using the `goopts` library. It showcases the ability to define multiple argument groups, each with its own set of arguments, allowing for a more organized and structured command-line interface.

## Program Description

### Arguments

- **Positional Arguments**:
  - `a (string)`: Help message for A.
  - `b (int)`: Help message for B.

- **String Arguments**:
  - `-c`, `--c (string)`: Help message for C.
  - `-d`, `--d (list of strings)`: Help message for D.
  - `-l`, `--l (string)`: Help message for L (required in the required mutually exclusive group).
  - `-m`, `--m (string)`: Help message for M (required in the required mutually exclusive group).
  - `-n`, `--n (string)`: Help message for N (not required in the not-required mutually exclusive group).
  - `-o`, `--o (string)`: Help message for O (not required in the not-required mutually exclusive group).
  - `-p`, `--p (string)`: Help message for P (required in the dependent argument group).
  - `-q`, `--q (string)`: Help message for Q (required in the dependent argument group).

- **Integer Arguments**:
  - `-e`, `--e (int)`: Help message for E.
  - `-f`, `--f (int range)`: Help message for F (value must be between 1 and 10).
  - `-g`, `--g (list of ints)`: Help message for G.

- **Map Arguments**:
  - `-i`, `--i (map of string to string)`: Help message for I (HTTP headers).

- **Boolean Arguments**:
  - `-j`, `--j (boolean)`: Help message for J.

- **TCP Port Argument**:
  - `-k`, `--k (int)`: Help message for K.

### How to Run

Compile and run the example program using the following commands:

```bash
go build -o poc
./poc "a" 17 -c valueC -l valueL
```

### Example Execution

- Execute with required and optional arguments:
    ```
    ./poc "a" 17 -c valueC -d "item1" -d "item1" -d "item3" -e 10 -f 5 -g 2 -g 3 -g 4 -i "Header1:Value1" -i "Header2:Value2" -j -k 8080 -l valueL
    ```
    Output:
    ```
        a: a
        b: 17
        c: valueC
        d: [item1 item1 item3]
        e: 10
        f: 5
        g: [2 3 4]
        i: map[Header1:Value1 Header2:Value2]
        j: true
        k: 8080
        l: valueL
        m: 
        n: 
        o: 
        p: 
        q: 
    ```

- Execute with all required arguments:
    ```
    ./poc "a" 17 -c valueC -l valueL
    ```
    Output:
    ```
        a: a
        b: 17
        c: valueC
        d: []
        e: 0
        f: 0
        g: []
        i: map[]
        j: false
        k: 0
        l: valueL
        m: 
        n: 
        o: 
        p: 
        q:
    ```

### Code Walkthrough

The program starts by setting up an `ArgumentsParser` with various argument types and groups:

```go
package main

import (
	"fmt"

	"github.com/p0dalirius/goopts/parser"
)

var (
	a string
	b int
	c string
	d []string
	e int
	f int
	g []int
	i map[string]string
	j bool
	k int
	l string
	m string
	n string
	o string
	p string
	q string
)

func parseArgs() {
	// Create a new arguments parser with a custom banner
	ap := parser.ArgumentsParser{Banner: "PoC of goopts parsing v.1.1 - by Remi GASCOU (Podalirius)"}

	ap.NewStringPositionalArgument(&a, "a", "Help message for A.")
	ap.NewIntPositionalArgument(&b, "b", "Help message for B.")

	ap.NewStringArgument(&c, "-c", "--c", "", true, "Help message for C.")
	ap.NewListOfStringsArgument(&d, "-d", "--d", []string{}, true, "Help message for D.")
	ap.NewIntArgument(&e, "-e", "--e", 0, true, "Help message for E.")
	ap.NewIntRangeArgument(&f, "-f", "--f", 5, 1, 10, true, "Help message for F.")
	ap.NewListOfIntsArgument(&g, "-g", "--g", []int{}, true, "Help message for G.")
	ap.NewMapOfHttpHeadersArgument(&i, "-i", "--i", map[string]string{}, true, "Help message for I.")
	ap.NewBoolArgument(&j, "-j", "--j", false, "Help message for J.")
	ap.NewTcpPortArgument(&k, "-k", "--k", 0, true, "Help message for K.");

	// Define an argument group for required mutually exclusive arguments
	group_rmeag, err := ap.NewRequiredMutuallyExclusiveArgumentGroup("1. RequiredMutuallyExclusiveArgumentGroup")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_rmeag.NewStringArgument(&l, "-l", "--l", "", true, "Help message for L.")
		group_rmeag.NewStringArgument(&m, "-m", "--m", "", true, "Help message for M.")
	}

	// Define an argument group for not-required mutually exclusive arguments
	group_nrmeag, err := ap.NewNotRequiredMutuallyExclusiveArgumentGroup("2. NotRequiredMutuallyExclusiveArgumentGroup")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_nrmeag.NewStringArgument(&n, "-n", "--n", "", true, "Help message for N.")
		group_nrmeag.NewStringArgument(&o, "-o", "--o", "", true, "Help message for O.")
	}

	// Define an argument group for dependent arguments
	group_dag, err := ap.NewDependentArgumentGroup("3. DependentArgumentGroup")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_dag.NewStringArgument(&p, "-p", "--p", "", true, "Help message for P.")
		group_dag.NewStringArgument(&q, "-q", "--q", "", true, "Help message for Q.")
	}

	// Parse the flags
	ap.Parse()
}

func main() {
	parseArgs()

	fmt.Printf("a: %v\n", a)
	fmt.Printf("b: %v\n", b)
	fmt.Printf("c: %v\n", c)
	fmt.Printf("d: %v\n", d)
	fmt.Printf("e: %v\n", e)
	fmt.Printf("f: %v\n", f)
	fmt.Printf("g: %v\n", g)
	fmt.Printf("i: %v\n", i)
	fmt.Printf("j: %v\n", j)
	fmt.Printf("k: %v\n", k)
	fmt.Printf("l: %v\n", l)
	fmt.Printf("m: %v\n", m)
	fmt.Printf("n: %v\n", n)
	fmt.Printf("o: %v\n", o)
	fmt.Printf("p: %v\n", p)
	fmt.Printf("q: %v\n", q)
}
```

## Explanation

- Positional Arguments:
    - `a` and `b` are defined as positional arguments with their respective help messages.

- String Arguments:
    - The string arguments `c`, `d`, and others provide a way to capture user input with appropriate descriptions.

- Integer and List Arguments:
    - The program demonstrates how to handle integers and lists, including ranges for the -f argument.

- Mutually Exclusive Groups:
    - Required and not-required mutually exclusive argument groups are defined to manage which arguments can be specified together.

- Dependent Arguments:
    - The dependent argument group shows how to require specific arguments based on the presence of others.

This structured approach to defining command-line arguments makes the application more user-friendly and maintains clarity in the functionality offered.
