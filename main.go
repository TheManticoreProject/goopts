package main

import (
	"fmt"

	"github.com/p0dalirius/goopts/subparser"
)

var (
	mode string

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
	asp := subparser.ArgumentsSubparser{
		Banner: "PoC of goopts parsing v.1.1 - by Remi GASCOU (Podalirius)",
		Name:   "mode",
		Value:  &mode,
	}

	// Define positional subparsers
	subparser_recon := asp.AddSubParser("recon", "Reconnaissance mode.")
	subparser_recon.NewStringArgument(&filePath, "f", "file", "The file to recon.", true, "The file to recon.")
	subparser_recon.NewStringArgument(&outputFolder, "o", "output", "The output folder.", true, "The output folder.")
	subparser_recon.NewBoolArgument(&enableLogging, "-l", "--enable-logging", true, "Enable logging during execution.")
	subparser_recon.NewBoolArgument(&disableEncryption, "-e", "--disable-encryption", false, "Disable encryption for data transfer.")
	// Define an argument group for database authentication
	group_dbAuth, err := subparser_recon.NewArgumentGroup("Database Authentication")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_dbAuth.NewStringArgument(&dbHost, "-H", "--db-host", "", true, "Hostname or IP of the database server.")
		group_dbAuth.NewStringArgument(&dbUsername, "-U", "--db-username", "", true, "Username for database authentication.")
		group_dbAuth.NewStringArgument(&dbPassword, "-P", "--db-password", "", true, "Password for database authentication.")
		group_dbAuth.NewTcpPortArgument(&dbPort, "-p", "--db-port", 3306, false, "Port number of the database server.")
	}

	subparser_scan := asp.AddSubParser("scan", "Scan mode.")
	subparser_scan.NewStringArgument(&filePath, "f", "file", "The file to scan.", true, "The file to scan.")
	// Define an argument group for server settings
	group_serverSettings, err := subparser_scan.NewArgumentGroup("Server Settings")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_serverSettings.NewStringArgument(&serverIP, "-i", "--server-ip", "", true, "IP address of the server to connect.")
		group_serverSettings.NewTcpPortArgument(&serverPort, "-s", "--server-port", 8080, false, "Port on which the server listens.")
	}

	// Parse the flags
	asp.Parse()
}

func main() {
	parseArgs()

	fmt.Printf("mode: %s\n\n", mode)

	fmt.Printf("[+] recon\n")
	fmt.Printf("  | filePath: %s\n", filePath)
	fmt.Printf("  | outputFolder: %s\n", outputFolder)
	fmt.Printf("  | enableLogging: %t\n", enableLogging)
	fmt.Printf("  | disableEncryption: %t\n", disableEncryption)
	fmt.Printf("  | dbHost: %s\n", dbHost)
	fmt.Printf("  | dbUsername: %s\n", dbUsername)
	fmt.Printf("  | dbPassword: %s\n", dbPassword)
	fmt.Printf("  | dbPort: %d\n", dbPort)
	fmt.Printf("  | serverIP: %s\n", serverIP)
	fmt.Printf("  | serverPort: %d\n\n", serverPort)

	fmt.Printf("[+] scan\n")
	fmt.Printf("  | filePath: %s\n", filePath)
	fmt.Printf("  | serverIP: %s\n", serverIP)
	fmt.Printf("  | serverPort: %d\n", serverPort)
}
