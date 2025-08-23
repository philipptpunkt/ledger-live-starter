package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Run setup to configure ledger-live-starter",
	Long:  `Initialize ledger-live-starter configuration with your preferences.`,
	Run:   runSetupCmd,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetupCmd(cmd *cobra.Command, args []string) {
	fmt.Printf("%s Welcome to Ledger Live Starter setup!\n\n", BoldColorText("Setup Mode", Cyan))
	
	_, err := runSetupMode()
	if err != nil {
		fmt.Printf("%s Setup failed: %v\n", ColorText("Error:", Red), err)
		return
	}
}

func runSetupMode() (*Config, error) {
	config := getDefaultConfig()
	
	fmt.Printf("%s Setting up configuration at: %s\n\n", BoldColorText("Info:", Blue), getConfigPath())

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

	err := pathForm.Run()
	if err != nil {
		return nil, fmt.Errorf("setup cancelled: %v", err)
	}

	config.LedgerLivePath = strings.TrimSpace(ledgerLivePath)

	// Step 2: Show default parameters and ask if user wants to modify
	fmt.Printf("\n%s Default parameters available:\n", BoldColorText("Parameters:", Yellow))
	for i, param := range config.Parameters {
		fmt.Printf("  %d. %s - %s\n", i+1, BoldText(param.Name), param.Description)
		fmt.Printf("     %s\n", ColorText(param.EnvVar, Cyan))
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

	err = addParamsForm.Run()
	if err != nil {
		return nil, fmt.Errorf("setup cancelled: %v", err)
	}

	// Step 3: Add custom parameters if requested
	if addMoreParams {
		config = addCustomParameters(config)
	}

	// Step 4: Save configuration
	err = ensureConfigDirExists()
	if err != nil {
		return nil, fmt.Errorf("failed to create config directory: %v", err)
	}

	err = saveConfigToPath(config, getConfigPath())
	if err != nil {
		return nil, fmt.Errorf("failed to save configuration: %v", err)
	}

	// Step 5: Show completion message
	fmt.Printf("\n%s Setup completed successfully!\n", ColorText("Success:", Green))
	fmt.Printf("Configuration saved to: %s\n\n", BoldText(getConfigPath()))
	
	fmt.Printf("%s How to use:\n", BoldColorText("Next Steps:", Cyan))
	fmt.Printf("  %s                    - Start with interactive menu\n", BoldText("ledger-live start"))
	fmt.Printf("  %s --config /path     - Use custom config file\n", BoldText("ledger-live start"))
	fmt.Printf("  %s                    - Run setup again\n", BoldText("ledger-live setup"))
	fmt.Println()

	return config, nil
}

func addCustomParameters(config *Config) *Config {
	fmt.Printf("\n%s Add custom parameters:\n", BoldColorText("Custom Setup:", Green))
	
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
		
		if err := nameForm.Run(); err != nil {
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
		
		if err := envForm.Run(); err != nil {
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
		
		descForm.Run() // Optional field, ignore errors

		// Add parameter
		newParam := Parameter{
			Name:        strings.TrimSpace(name),
			EnvVar:      strings.TrimSpace(envVar),
			Description: strings.TrimSpace(description),
		}
		config.Parameters = append(config.Parameters, newParam)
		
		fmt.Printf("%s Added parameter: %s\n", ColorText("✓", Green), newParam.Name)

		// Ask for another
		anotherForm := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Add another parameter?").
					Value(&addAnother),
			),
		)
		
		if err := anotherForm.Run(); err != nil || !addAnother {
			break
		}
	}
	
	return config
}
