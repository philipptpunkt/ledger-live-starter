package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"ledger-live-starter/cmd/ledger-live/presets"
	"ledger-live-starter/cmd/ledger-live/setup"
	"ledger-live-starter/cmd/ledger-live/ui"
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
	fmt.Println(ui.GetLogo())
	fmt.Println()
	fmt.Println(getVersionOrUpdateDisplay())
	fmt.Println()
	
	// Load configuration
	config, err := setup.LoadConfig()
	if err != nil {
		fmt.Printf("%s %s\n\n", WarningText("Warning:"), NormalText(fmt.Sprintf("Could not load config.json (%v), using defaults", err)))
		config = setup.GetDefaultConfig()
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
	
	// Add presets with default styling
	for _, preset := range config.Presets {
		options = append(options, huh.NewOption(preset.Name, preset.Name))
	}
	
	// Add standard options
	options = append(options, huh.NewOption("Start manually", "manual"))
	options = append(options, huh.NewOption("More", "more"))
	options = append(options, huh.NewOption("Exit", "exit"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose an option:").
				Description("Select a preset to start the application directly."). // make dynamic based on selected option
				Options(options...).
				Value(&selected),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		ShowCancellationMessage()
		return
	}

	// Handle selection
	switch selected {
	case "manual":
		startManually()
	case "more":
		showMoreMenu(config)
	case "exit":
		ShowGoodbyeMessage()
		return
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
	options = append(options, huh.NewOption("Exit", "exit"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("No presets found. Choose an option:").
				Options(options...).
				Value(&selected),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		ShowCancellationMessage()
		return
	}

	switch selected {
	case "create":
		presets.CreatePreset()
	case "manual":
		startManually()
	case "more":
		showMoreMenu(nil) // Pass nil since no presets exist
	case "exit":
		ShowGoodbyeMessage()
		return
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
		fmt.Printf("%s %s '%s' %s\n", ErrorText("❌"), NormalText("Preset"), HighlightText(presetName), NormalText("not found"))
		return
	}

	// Convert preset to command
	cmdInfo := buildPresetCommand(selectedPreset, config)
	
	fmt.Printf("%s %s %s\n", SuccessText("✅"), NormalText("Starting preset:"), HighlightText(selectedPreset.Name))
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
	options = append(options, huh.NewOption("Edit presets", "edit"))
	
	// Always show "Edit parameters" option
	options = append(options, huh.NewOption("Edit parameters", "parameters"))
	
	// Always show back option
	options = append(options, huh.NewOption("Back", "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("More options:").
				Options(options...).
				Value(&selected),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		ShowCancellationMessage()
		return
	}

	switch selected {
	case "edit":
		presets.EditPresets()
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

// Deprecated: Use getVersionOrUpdateDisplay() instead
// Kept for backward compatibility, but now just shows current version
func getVersionDisplay() string {
	// Define adaptive text color (same as theme.go)
	textColor := lipgloss.AdaptiveColor{
		Light: "#1a1a1a", // Dark text for light backgrounds
		Dark:  "#ffffff", // White text for dark backgrounds
	}
	
	// Create version text with adaptive color and center alignment
	versionStyle := lipgloss.NewStyle().
		Foreground(textColor).
		Align(lipgloss.Center).
		MarginLeft(4) // Same left margin as logo
	
	return versionStyle.Render(fmt.Sprintf("v%s", Version))
}