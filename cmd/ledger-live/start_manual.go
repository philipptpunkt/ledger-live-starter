package main

import (
	"fmt"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// Manual start flow
func startManually() {
	// Load configuration
	config, err := setup.LoadConfig()
	if err != nil {
		fmt.Printf("%s %s\n\n", WarningText("Warning:"), NormalText(fmt.Sprintf("Could not load config.json (%v), using defaults", err)))
		config = setup.GetDefaultConfig()
	}

	// Step 1: Platform selection
	platform, baseCommand, err := selectPlatform()
	if err != nil {
		fmt.Printf("%s %v\n", ErrorText("Error:"), NormalText(err.Error()))
		return
	}

	// Step 2: Parameter selection
	selectedParams, err := selectParameters(config.Parameters)
	if err != nil {
		fmt.Printf("%s %v\n", ErrorText("Error:"), NormalText(err.Error()))
		return
	}

	// Step 3: Build and execute command
	cmdInfo := buildCommand(baseCommand, selectedParams, config)
	fmt.Printf("\n%s %s %s...\n", SuccessText("Success:"), NormalText("Starting"), HighlightText(platform))
	executeCommand(cmdInfo)
}
