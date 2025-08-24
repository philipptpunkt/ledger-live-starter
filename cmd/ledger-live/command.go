package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

)

type CommandInfo struct {
	BaseCommand    string
	EnvVars        map[string]string
	WorkingDir     string
}

func buildCommand(baseCommand string, parameters []Parameter, config *Config) *CommandInfo {
	envVars := make(map[string]string)

	// Extract environment variables from selected parameters
	for _, param := range parameters {
		// Parse "VAR_NAME=value" format
		if strings.Contains(param.EnvVar, "=") {
			parts := strings.SplitN(param.EnvVar, "=", 2)
			if len(parts) == 2 {
				envVars[parts[0]] = parts[1]
			}
		}
	}

	return &CommandInfo{
		BaseCommand: baseCommand,
		EnvVars:     envVars,
		WorkingDir:  config.LedgerLivePath,
	}
}

func executeCommand(cmdInfo *CommandInfo) {
	// Build display string for user
	var displayParts []string
	for key, value := range cmdInfo.EnvVars {
		displayParts = append(displayParts, fmt.Sprintf("%s=%s", key, value))
	}
	
	displayCommand := strings.Join(displayParts, " ")
	if len(displayParts) > 0 {
		displayCommand += " " + cmdInfo.BaseCommand
	} else {
		displayCommand = cmdInfo.BaseCommand
	}
	
	fmt.Printf("\n%s %s\n", TitleText("Executing:"), HighlightText(displayCommand))
	fmt.Printf("%s %s\n\n", InfoTextTitle("Working directory:"), HighlightText(cmdInfo.WorkingDir))

	// Split base command into parts
	parts := strings.Fields(cmdInfo.BaseCommand)
	if len(parts) == 0 {
		fmt.Printf("%s %s\n", ErrorText("Error:"), NormalText("Invalid command"))
		return
	}

	// Create command
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	// Set working directory
	if cmdInfo.WorkingDir != "" {
		cmd.Dir = cmdInfo.WorkingDir
	}
	
	// Set environment variables
	cmd.Env = os.Environ() // Start with current environment
	for key, value := range cmdInfo.EnvVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// Execute command
	err := cmd.Run()
	if err != nil {
		fmt.Printf("%s %s %s\n", ErrorText("Error:"), NormalText("Error executing command:"), NormalText(err.Error()))
		os.Exit(1)
	}
}
