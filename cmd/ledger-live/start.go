package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the ledger-live application",
	Long:  `Start ledger-live using presets or manual configuration.`,
	Run:   runStartCmd,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func runStartCmd(cmd *cobra.Command, args []string) {
	fmt.Println(BoldColorText("Ledger Live Starter", Cyan))
	fmt.Println()

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("%s Warning: Could not load config.json (%v), using defaults\n\n", BoldColorText("!", Yellow), err)
		config = getDefaultConfig()
	}

	// Show main menu based on preset availability
	if len(config.Presets) > 0 {
		showPresetMenu(config)
	} else {
		showNoPresetMenu()
	}
}

func showPresetMenu(config *Config) {
	var options []huh.Option[string]
	
	// Add presets
	for _, preset := range config.Presets {
		options = append(options, huh.NewOption(preset.Name, preset.Name))
	}
	
	// Add standard options
	options = append(options, huh.NewOption("Start manually", "manual"))
	options = append(options, huh.NewOption("More", "more"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose an option:").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("\n❌ Selection cancelled")
		return
	}

	// Handle selection
	switch selected {
	case "manual":
		startManually()
	case "more":
		showMoreMenu(config)
	default:
		// User selected a preset
		executePreset(selected, config)
	}
}

func showNoPresetMenu() {
	var options []huh.Option[string]
	options = append(options, huh.NewOption("Create preset", "create"))
	options = append(options, huh.NewOption("Start manually", "manual"))
	options = append(options, huh.NewOption("More", "more"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("No presets found. Choose an option:").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("\n❌ Selection cancelled")
		return
	}

	switch selected {
	case "create":
		createPreset()
	case "manual":
		startManually()
	case "more":
		showMoreMenu(nil) // Pass nil since no presets exist
	}
}

func executePreset(presetName string, config *Config) {
	// Find the preset
	var selectedPreset *Preset
	for _, preset := range config.Presets {
		if preset.Name == presetName {
			selectedPreset = &preset
			break
		}
	}

	if selectedPreset == nil {
		fmt.Printf("❌ Preset '%s' not found\n", presetName)
		return
	}

	// Convert preset to command
	cmdInfo := buildPresetCommand(selectedPreset, config)
	
	fmt.Printf("✅ Starting preset: %s\n", selectedPreset.Name)
	executeCommand(cmdInfo)
}

func buildPresetCommand(preset *Preset, config *Config) *CommandInfo {
	// Determine base command from platform
	var baseCommand string
	switch preset.Platform {
	case "mobile":
		baseCommand = "pnpm dev:llm"
	case "desktop":
		baseCommand = "pnpm dev:lld"
	default:
		baseCommand = "pnpm dev:llm" // default to mobile
	}

	// Find and parse parameter env vars
	envVars := make(map[string]string)
	for _, paramName := range preset.Parameters {
		for _, param := range config.Parameters {
			if param.Name == paramName {
				// Parse "VAR_NAME=value" format
				if strings.Contains(param.EnvVar, "=") {
					parts := strings.SplitN(param.EnvVar, "=", 2)
					if len(parts) == 2 {
						envVars[parts[0]] = parts[1]
					}
				}
				break
			}
		}
	}

	return &CommandInfo{
		BaseCommand: baseCommand,
		EnvVars:     envVars,
		WorkingDir:  config.LedgerLivePath,
	}
}

func showMoreMenu(config *Config) {
	var options []huh.Option[string]
	
	// Always show "Edit presets" option
	options = append(options, huh.NewOption(BoldColorText("Edit presets", Green), "edit"))
	
	// Always show "Edit parameters" option
	options = append(options, huh.NewOption(BoldColorText("Edit parameters", Yellow), "parameters"))
	
	// Always show back option
	options = append(options, huh.NewOption(ColorText("Back", Red), "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("More options:").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("\n❌ Selection cancelled")
		return
	}

	switch selected {
	case "edit":
		editPresets()
	case "parameters":
		editParameters()
	case "back":
		// Go back to the appropriate menu
		if config != nil && len(config.Presets) > 0 {
			showPresetMenu(config)
		} else {
			showNoPresetMenu()
		}
	}
}