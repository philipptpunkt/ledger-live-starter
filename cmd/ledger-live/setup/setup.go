package setup

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var SetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run setup to configure ledger-live-starter",
	Long:  `Initialize ledger-live-starter configuration with your preferences.`,
	Run:   runSetupCmd,
}

func runSetupCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("%s %s\n\n", TitleText("Setup Mode"), NormalText("Welcome to Ledger Live Starter setup!"))
	
	_, err := RunSetupMode()
	if err != nil {
		fmt.Printf("%s %s %s\n", ErrorText("Error:"), NormalText("Setup failed:"), NormalText(err.Error()))
		return
	}
}

func RunSetupMode() (*Config, error) {
	config := GetDefaultConfig()
	
	fmt.Printf("%s %s %s\n\n", TitleText("Info:"), NormalText("Setting up configuration at:"), HighlightText(GetConfigPath()))

	// Step 1: Get Ledger Live path
	var ledgerLivePath string
	pathForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter the absolute path to your ledger-live directory:").
				Description("This is where you cloned the ledger-live repository").
				Placeholder("e.g., /Users/yourname/projects/ledger-live").
				Value(&ledgerLivePath).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("ledger-live path cannot be empty")
					}
					if !filepath.IsAbs(s) {
						return fmt.Errorf("path must be absolute (start with / or C:\\)")
					}
					return nil
				}),
		),
	)

	err := RunStyledForm(pathForm)
	if err != nil {
		return nil, fmt.Errorf("setup cancelled: %v", err)
	}

	config.LedgerLivePath = strings.TrimSpace(ledgerLivePath)

	// Step 2: Show default parameters and ask if user wants to modify
	fmt.Printf("\n%s %s\n", TitleText("Parameters:"), NormalText("Default parameters available:"))
	for i, param := range config.Parameters {
		fmt.Printf("  %d. %s - %s\n", i+1, HighlightText(param.Name), NormalText(param.Description))
		fmt.Printf("     %s\n", NormalText(param.EnvVar))
	}
	fmt.Println()

	var addMoreParams bool
	addParamsForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Would you like to add additional parameters?").
				Value(&addMoreParams),
		),
	)

	err = RunStyledForm(addParamsForm)
	if err != nil {
		return nil, fmt.Errorf("setup cancelled: %v", err)
	}

	// Step 3: Add custom parameters if requested
	if addMoreParams {
		config = addCustomParameters(config)
	}

	// Step 4: Save configuration
	err = EnsureConfigDirExists()
	if err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	err = SaveConfigToPath(config, GetConfigPath())
	if err != nil {
		return nil, fmt.Errorf("failed to save configuration: %v", err)
	}

	// Step 5: Show completion message
	fmt.Printf("\n%s %s\n", SuccessText("Success:"), NormalText("Setup completed successfully!"))
	fmt.Printf("%s %s\n\n", NormalText("Configuration saved to:"), HighlightText(GetConfigPath()))
	
	fmt.Printf("%s %s\n", TitleText("Next Steps:"), NormalText("How to use:"))
	fmt.Printf("  %s %s\n", HighlightText("ledger-live start"), NormalText("                    - Start with interactive menu"))
	fmt.Printf("  %s %s\n", HighlightText("ledger-live start --config /path"), NormalText("     - Use custom config file"))
	fmt.Printf("  %s %s\n", HighlightText("ledger-live setup"), NormalText("                    - Run setup again"))
	fmt.Println()

	return config, nil
}

func addCustomParameters(config *Config) *Config {
	fmt.Printf("\n%s %s\n", TitleText("Custom Setup:"), NormalText("Add custom parameters:"))
	
	for {
		var addAnother bool
		var name, envVar, description string

		// Parameter name
		nameForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Parameter name:").
					Placeholder("e.g., 'Enable debug logs'").
					Value(&name),
			),
		)
		
		if err := RunStyledForm(nameForm); err != nil {
			break
		}
		
		if strings.TrimSpace(name) == "" {
			break
		}

		// Environment variable
		envForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Environment variable:").
					Placeholder("e.g., 'DEBUG_LOGS=1'").
					Value(&envVar).
					Validate(func(s string) error {
						if !strings.Contains(s, "=") {
							return fmt.Errorf("must include '=' (e.g., VAR_NAME=value)")
						}
						return nil
					}),
			),
		)
		
		if err := RunStyledForm(envForm); err != nil {
			break
		}

		// Description
		descForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Description (optional):").
					Placeholder("e.g., 'Enable detailed debug logging'").
					Value(&description),
			),
		)
		
		RunStyledForm(descForm) // Optional field, ignore errors

		// Add parameter
		newParam := Parameter{
			Name:        strings.TrimSpace(name),
			EnvVar:      strings.TrimSpace(envVar),
			Description: strings.TrimSpace(description),
		}
		config.Parameters = append(config.Parameters, newParam)
		
		fmt.Printf("%s %s %s\n", SuccessText("✓"), NormalText("Added parameter:"), HighlightText(newParam.Name))

		// Ask for another
		anotherForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another parameter?").
					Value(&addAnother),
			),
		)
		
		if err := RunStyledForm(anotherForm); err != nil || !addAnother {
			break
		}
	}
	
	return config
}

// These functions need to be imported from the main package
// We'll define them as variables that get set from main
var (
	RunStyledForm func(*huh.Form) error
	
	// Theme-based text functions
	TitleText     func(text string) string
	ErrorText     func(text string) string
	SuccessText   func(text string) string
	InfoTextTitle func(text string) string
	HighlightText func(text string) string
	NormalText    func(text string) string
	WarningText   func(text string) string
)
