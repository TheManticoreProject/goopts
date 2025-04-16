package main

import (
	"fmt"

	"github.com/TheManticoreProject/goopts/parser"
)

var (
	mode        string
	groupA_mode string

	// Positional arguments
	filePath          string
	outputFolder      string
	enableLogging     bool
	disableEncryption bool

	dbHost     string
	dbUsername string
	dbPassword string
	dbPort     int

	serverIP   string
	serverPort int
)

func parseArgs() {
	// Create a new arguments parser with a custom banner
	ap := parser.NewParser("PoC of goopts parsing v1.2.0 - by Remi GASCOU (Podalirius) @ TheManticoreProject")
	ap.SetupSubParsing("mode", &mode, true)

	// Define positional subparsers
	subparser_groupA := ap.AddSubParser("groupA", "groupA mode.")
	subparser_groupA.SetupSubParsing("groupA_mode", &groupA_mode, true)

	subparser_groupA_groupAB := subparser_groupA.AddSubParser("groupAB", "groupAB mode.")
	subparser_groupA_groupAB.NewBoolArgument(&enableLogging, "", "--enable-logging", true, "Enable logging during execution.")
	subparser_groupA_groupAB_server, err := subparser_groupA_groupAB.NewArgumentGroup("Server")
	if err != nil {
		fmt.Printf("[-] Error creating group: %s\n", err)
	} else {
		subparser_groupA_groupAB_server.NewStringArgument(&dbHost, "", "--db-host", "The database host.", true, "The database host.")
		subparser_groupA_groupAB_server.NewIntArgument(&serverPort, "", "--server-port", 1337, true, "The server port.")
	}

	subparser_groupA_groupAC := subparser_groupA.AddSubParser("groupAC", "groupAC mode.")
	subparser_groupA_groupAC.NewStringArgument(&dbHost, "", "--db-host", "The database host.", true, "The database host.")
	subparser_groupA_groupAC_server, err := subparser_groupA_groupAC.NewArgumentGroup("Server")
	if err != nil {
		fmt.Printf("[-] Error creating group: %s\n", err)
	} else {
		subparser_groupA_groupAC_server.NewStringArgument(&serverIP, "", "--server-ip", "The server IP.", true, "The server IP.")
		subparser_groupA_groupAC_server.NewIntArgument(&serverPort, "", "--server-port", 1337, true, "The server port.")
	}

	// Define positional subparsers
	subparser_groupZ := ap.AddSubParser("groupZ", "Add mode.")
	subparser_groupZ.NewStringArgument(&filePath, "f", "file", "The file to add.", true, "The file to add.")

	// Parse the flags
	ap.Parse()
}

func main() {
	parseArgs()

	fmt.Printf("mode: %s\n\n", mode)

	fmt.Printf("[+] groupA\n")
	fmt.Printf("  | Mode: %s\n", groupA_mode)
	fmt.Printf("  | Enable Logging: %t\n", enableLogging)
	fmt.Printf("  | Disable Encryption: %t\n", disableEncryption)
	fmt.Printf("  | Database Host: %s\n", dbHost)
	fmt.Printf("  | Database Username: %s\n", dbUsername)
	fmt.Printf("  | Database Password: %s\n", dbPassword)
	fmt.Printf("  | Database Port: %d\n", dbPort)
	fmt.Printf("  | Server IP: %s\n", serverIP)
	fmt.Printf("  | Server Port: %d\n\n", serverPort)

	fmt.Printf("[+] groupB\n")
	fmt.Printf("  | File Path: %s\n", filePath)
	fmt.Printf("  | Server IP: %s\n", serverIP)
	fmt.Printf("  | Server Port: %d\n", serverPort)
}
