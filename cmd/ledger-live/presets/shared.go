package presets

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// CommandInfo represents command information - duplicate here to avoid circular imports
type CommandInfo struct {
	BaseCommand string
	EnvVars     map[string]string
	WorkingDir  string
}

// loadConfigWithDefaults loads config with consistent error handling
func loadConfigWithDefaults() (*setup.Config, error) {
	config, err := setup.LoadConfig()
	if err != nil {
		fmt.Printf("%s %s\n\n", WarningText("Warning:"), NormalText(fmt.Sprintf("Could not load config.json (%v), using defaults", err)))
		return setup.GetDefaultConfig(), err
	}
	return config, nil
}

// convertPlatformToKey converts platform command to platform key
func convertPlatformToKey(platformCommand string) string {
	switch platformCommand {
	case "pnpm dev:llm":
		return "mobile"
	case "pnpm dev:lld":
		return "desktop"
	default:
		return "mobile"
	}
}

// extractParameterNames converts parameter structs to names
func extractParameterNames(selectedParams []setup.Parameter) []string {
	var parameterNames []string
	for _, param := range selectedParams {
		parameterNames = append(parameterNames, param.Name)
	}
	return parameterNames
}

// displayPresetSummary shows a formatted preset summary
func displayPresetSummary(preset setup.Preset) {
	fmt.Println(TitleText("Preset Summary:"))
	fmt.Printf("   %s %s\n", InfoTextTitle("Name:"), HighlightText(preset.Name))
	fmt.Printf("   %s %s\n", InfoTextTitle("Platform:"), HighlightText(strings.Title(preset.Platform)))
	if len(preset.Parameters) > 0 {
		fmt.Printf("   %s %s\n", InfoTextTitle("Parameters:"), NormalText(strings.Join(preset.Parameters, ", ")))
	} else {
		fmt.Printf("   %s %s\n", InfoTextTitle("Parameters:"), NormalText("None"))
	}
	fmt.Println()
}

// findPresetByName finds a preset by name and returns its index and pointer
func findPresetByName(presetName string, presets []setup.Preset) (int, *setup.Preset) {
	for i, preset := range presets {
		if preset.Name == presetName {
			return i, &presets[i]
		}
	}
	return -1, nil
}

// saveConfigWithError saves config and displays error if needed
func saveConfigWithError(config *setup.Config) error {
	err := setup.SaveConfig(config)
	if err != nil {
		fmt.Printf("%s %s %s\n", ErrorText("Error:"), NormalText("Error saving changes:"), NormalText(err.Error()))
	}
	return err
}

// Text styling function placeholders - these will be injected from main
var (
	TitleText     func(text string) string
	ErrorText     func(text string) string
	SuccessText   func(text string) string
	InfoTextTitle func(text string) string
	HighlightText func(text string) string
	NormalText    func(text string) string
	WarningText   func(text string) string
)

// Form execution placeholders - these will be injected from main
var (
	RunStyledForm              func(form *huh.Form) error
	ShowCancellationMessage    func()
	ShowConfirmationCancelledMessage func()
)

// UI function placeholders - these will be injected from main
var (
	InputPresetName    func(existingPresets []setup.Preset) (string, error)
	SelectPlatform     func() (string, string, error)
	SelectParameters   func(availableParams []setup.Parameter) ([]setup.Parameter, error)
	BuildPresetCommand func(preset *setup.Preset, config *setup.Config) *CommandInfo
	ExecuteCommand     func(cmdInfo *CommandInfo)
	ShowMoreMenu       func(config *setup.Config)
)
