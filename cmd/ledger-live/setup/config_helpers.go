package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
var configPath string

func SetConfigPath(path string) {
	configPath = path
}

func GetConfigPath() string {
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

func EnsureConfigDirExists() error {
	return os.MkdirAll(getLedgerLiveDir(), 0755)
}

func ConfigExists() bool {
	_, err := os.Stat(GetConfigPath())
	return err == nil
}

func SaveConfigToPath(config *Config, path string) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

func LoadConfig() (*Config, error) {
	// Check if config exists, trigger setup if not
	if !ConfigExists() {
		fmt.Printf("%s %s\n\n", TitleText("Setup:"), NormalText("No configuration found. Running setup mode..."))
		return RunSetupMode()
	}
	
	configFilePath := GetConfigPath()
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

func GetDefaultConfig() *Config {
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

func SaveConfig(config *Config) error {
	return SaveConfigToPath(config, GetConfigPath())
}
