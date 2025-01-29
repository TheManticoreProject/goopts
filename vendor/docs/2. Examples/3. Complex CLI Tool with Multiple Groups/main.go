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
	// h
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
	ap.NewTcpPortArgument(&k, "-k", "--k", 0, true, "Help message for K.")

	// Define an argument group for database authentication
	group_rmeag, err := ap.NewRequiredMutuallyExclusiveArgumentGroup("1. RequiredMutuallyExclusiveArgumentGroup")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_rmeag.NewStringArgument(&l, "-l", "--l", "", true, "Help message for L.")
		group_rmeag.NewStringArgument(&m, "-m", "--m", "", true, "Help message for M.")
	}

	// Define an argument group for database authentication
	group_nrmeag, err := ap.NewNotRequiredMutuallyExclusiveArgumentGroup("2. NotRequiredMutuallyExclusiveArgumentGroup")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_nrmeag.NewStringArgument(&n, "-n", "--n", "", true, "Help message for N.")
		group_nrmeag.NewStringArgument(&o, "-o", "--o", "", true, "Help message for O.")
	}

	// Define an argument group for database authentication
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
