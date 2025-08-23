package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config structures
type Config struct {
	LedgerLivePath string      `json:"ledger-live-path"`
	Parameters     []Parameter `json:"parameters"`
	Presets        []Preset    `json:"presets,omitempty"`
}

type Parameter struct {
	Name        string `json:"name"`
	EnvVar      string `json:"env_var"`
	Description string `json:"description"`
}

type Preset struct {
	Name       string   `json:"name"`
	Platform   string   `json:"platform"`    // "mobile" or "desktop"
	Parameters []string `json:"parameters"`  // List of parameter names
}

// Config path management functions

func getConfigPath() string {
	// Use flag value if provided
	if configPath != "" {
		return configPath
	}
	
	// Check environment variable
	if customPath := os.Getenv("LEDGER_LIVE_STARTER_CONFIG"); customPath != "" {
		return customPath
	}
	
	// Use ~/.ledger-live/config.json (next to binary)
	return getLedgerLiveConfigPath()
}

func getLedgerLiveConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		return "config.json"
	}
	
	return filepath.Join(homeDir, ".ledger-live", "config.json")
}

func getLedgerLiveDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	
	return filepath.Join(homeDir, ".ledger-live")
}

func ensureConfigDirExists() error {
	return os.MkdirAll(getLedgerLiveDir(), 0755)
}

func configExists() bool {
	_, err := os.Stat(getConfigPath())
	return err == nil
}

func saveConfigToPath(config *Config, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

// Manual start flow (moved from start.go)
func startManually() {
	fmt.Println(BoldColorText("Manual Start", Purple))
	fmt.Println()

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("%s Warning: Could not load config.json (%v), using defaults\n\n", BoldColorText("!", Yellow), err)
		config = getDefaultConfig()
	}

	// Step 1: Platform selection
	platform, baseCommand, err := selectPlatform()
	if err != nil {
		fmt.Printf("%s %v\n", ColorText("Error:", Red), err)
		return
	}

	// Step 2: Parameter selection
	selectedParams, err := selectParameters(config.Parameters)
	if err != nil {
		fmt.Printf("%s %v\n", ColorText("Error:", Red), err)
		return
	}

	// Step 3: Build and execute command
	cmdInfo := buildCommand(baseCommand, selectedParams, config)
	fmt.Printf("\n%s Starting %s...\n", ColorText("Success:", Green), platform)
	executeCommand(cmdInfo)
}

func loadConfig() (*Config, error) {
	// Check if config exists, trigger setup if not
	if !configExists() {
		fmt.Printf("%s No configuration found. Running setup mode...\n\n", BoldColorText("Setup:", Cyan))
		return runSetupMode()
	}
	
	configFilePath := getConfigPath()
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configFilePath, err)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

func getDefaultConfig() *Config {
	return &Config{
		LedgerLivePath: "",
		Parameters: []Parameter{
			{
				Name:        "Skip onboarding",
				EnvVar:      "SKIP_ONBOARDING=1",
				Description: "Enable skipping the onboarding process on mobile",
			},
			{
				Name:        "Disable transaction broadcast",
				EnvVar:      "DISABLE_TRANSACTION_BROADCAST=1",
				Description: "Disable broadcasting transactions and directly get success",
			},
			{
				Name:        "Bypass CORS",
				EnvVar:      "BYPASS_CORS=1",
				Description: "Bypass CORS restrictions for locale development",
			},
		},
		Presets: []Preset{},
	}
}

// Platform and parameter selection functions moved to shared.go

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
	
	fmt.Printf("\n%s %s\n", BoldColorText("Executing:", Blue), displayCommand)
	fmt.Printf("%s %s\n\n", ColorText("Working directory:", Cyan), cmdInfo.WorkingDir)

	// Split base command into parts
	parts := strings.Fields(cmdInfo.BaseCommand)
	if len(parts) == 0 {
		fmt.Printf("%s Invalid command\n", ColorText("Error:", Red))
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
		fmt.Printf("%s Error executing command: %v\n", ColorText("Error:", Red), err)
		os.Exit(1)
	}
}