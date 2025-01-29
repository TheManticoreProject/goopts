package main

import (
	"fmt"

	"github.com/p0dalirius/goopts/parser"
)

var (
	// Positional arguments
	filePath     string
	outputFolder string

	// Authentication flags
	dbHost     string
	dbUsername string
	dbPassword string
	dbPort     int

	// Server settings flags
	serverIP   string
	serverPort int

	// Feature flags
	enableLogging     bool
	disableEncryption bool
)

func parseArgs() {
	// Create a new arguments parser with a custom banner
	ap := parser.ArgumentsParser{Banner: "PoC of goopts parsing v.1.1 - by Remi GASCOU (Podalirius)"}

	// Define positional arguments
	ap.NewStringPositionalArgument(&filePath, "filepath", "Path to the input file.")
	ap.NewStringPositionalArgument(&outputFolder, "outputfolder", "Destination folder for output.")

	// Define global flags
	ap.NewBoolArgument(&enableLogging, "-l", "--enable-logging", true, "Enable logging during execution.")
	ap.NewBoolArgument(&disableEncryption, "-e", "--disable-encryption", false, "Disable encryption for data transfer.")

	// Define an argument group for database authentication
	group_dbAuth, err := ap.NewArgumentGroup("Database Authentication")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_dbAuth.NewStringArgument(&dbHost, "-H", "--db-host", "", true, "Hostname or IP of the database server.")
		group_dbAuth.NewStringArgument(&dbUsername, "-U", "--db-username", "", true, "Username for database authentication.")
		group_dbAuth.NewStringArgument(&dbPassword, "-P", "--db-password", "", true, "Password for database authentication.")
		group_dbAuth.NewTcpPortArgument(&dbPort, "-p", "--db-port", 3306, false, "Port number of the database server.")
	}

	// Define an argument group for server settings
	group_serverSettings, err := ap.NewArgumentGroup("Server Settings")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_serverSettings.NewStringArgument(&serverIP, "-i", "--server-ip", "", true, "IP address of the server to connect.")
		group_serverSettings.NewTcpPortArgument(&serverPort, "-s", "--server-port", 8080, false, "Port on which the server listens.")
	}

	// Parse the flags
	ap.Parse()
}

func main() {
	parseArgs()
}
