package main

import (
	"fmt"

	"github.com/TheManticoreProject/goopts/subparser"
)

var (
	// Positional arguments
	filePath     string
	outputFolder string
	timeout      int
	verbose      bool
	scanDepth    int
	quickScan    bool
)

func parseArgs() {
	// Create a new arguments parser with a custom banner
	asp := subparser.ArgumentsSubparser{
		Banner: "PoC of goopts parsing v.1.1 - by Remi GASCOU (Podalirius)",
		Name:   "mode",
	}

	// Define positional subparsers
	subparser_recon := asp.AddSubParser("recon", "Reconnaissance mode.")
	subparser_recon.NewStringArgument(&filePath, "f", "file", "The file to recon.", true, "The file to recon.")
	subparser_recon.NewStringArgument(&outputFolder, "o", "output", "The output folder.", true, "The output folder.")
	subparser_recon.NewIntArgument(&timeout, "t", "timeout", 30, true, "Timeout for the recon operation.")
	subparser_recon.NewBoolArgument(&verbose, "v", "verbose", false, "Enable verbose output for recon.")

	subparser_scan := asp.AddSubParser("scan", "Scan mode.")
	subparser_scan.NewStringArgument(&filePath, "f", "file", "The file to scan.", true, "The file to scan.")
	subparser_scan.NewIntArgument(&scanDepth, "d", "depth", 5, true, "Depth of the scan.")
	subparser_scan.NewBoolArgument(&quickScan, "q", "quick", false, "Enable quick scan mode.")

	// Parse the flags
	asp.Parse()
}

func main() {
	parseArgs()

	fmt.Printf("filePath: %s\n", filePath)
	fmt.Printf("outputFolder: %s\n", outputFolder)
	fmt.Printf("timeout: %d\n", timeout)
	fmt.Printf("verbose: %t\n", verbose)
	fmt.Printf("scanDepth: %d\n", scanDepth)
	fmt.Printf("quickScan: %t\n", quickScan)
}
