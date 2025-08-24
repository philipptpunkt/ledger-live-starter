package parameters

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// loadConfigWithDefaults loads config with consistent error handling
func loadConfigWithDefaults() (*setup.Config, error) {
	config, err := setup.LoadConfig()
	if err != nil {
		fmt.Printf("%s %s\n\n", WarningText("Warning:"), NormalText(fmt.Sprintf("Could not load config.json (%v), using defaults", err)))
		return setup.GetDefaultConfig(), err
	}
	return config, nil
}

// saveConfigWithError saves config and displays error if needed
func saveConfigWithError(config *setup.Config) error {
	err := setup.SaveConfig(config)
	if err != nil {
		fmt.Printf("%s %s %s\n", ErrorText("Error:"), NormalText("Error saving changes:"), NormalText(err.Error()))
	}
	return err
}

// findParameterByName finds a parameter by name and returns its index and pointer
func findParameterByName(paramName string, parameters []setup.Parameter) (int, *setup.Parameter) {
	for i, param := range parameters {
		if param.Name == paramName {
			return i, &parameters[i]
		}
	}
	return -1, nil
}

// validateParameterName validates parameter name for uniqueness and emptiness
func validateParameterName(name string, existingParams []setup.Parameter, currentParamName string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("parameter name cannot be empty")
	}
	
	for _, p := range existingParams {
		// Allow keeping the same name if we're editing
		if p.Name == name && p.Name != currentParamName {
			return fmt.Errorf("parameter '%s' already exists", name)
		}
	}
	return nil
}

// validateEnvironmentVariable validates environment variable format
func validateEnvironmentVariable(envVar string) error {
	if strings.TrimSpace(envVar) == "" {
		return fmt.Errorf("environment variable cannot be empty")
	}
	if !strings.Contains(envVar, "=") {
		return fmt.Errorf("environment variable must include '=' (e.g., VAR_NAME=value)")
	}
	return nil
}

// createParameterOptionsFromList creates huh options from parameter list
func createParameterOptionsFromList(parameters []setup.Parameter) []huh.Option[string] {
	var options []huh.Option[string]
	for _, param := range parameters {
		options = append(options, huh.NewOption(param.Name, param.Name))
	}
	return options
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
	RunStyledForm                    func(form *huh.Form) error
	ShowCancellationMessage          func()
	ShowConfirmationCancelledMessage func()
)

// Navigation function placeholder - this will be injected from main
var (
	ShowMoreMenu func(config *setup.Config)
)
