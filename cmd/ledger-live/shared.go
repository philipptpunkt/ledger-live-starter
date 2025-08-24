package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"

)

// Shared UI functions that can be reused across different flows

func selectPlatform() (string, string, error) {
	var options []huh.Option[string]
	options = append(options, huh.NewOption("Mobile", "mobile"))
	options = append(options, huh.NewOption("Desktop", "desktop"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Start Ledger-Live for:").
				Description("Select the platform you want to start Ledger-Live for.").
				Options(options...).
				Value(&selected),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		return "", "", fmt.Errorf("platform selection cancelled")
	}

	switch selected {
	case "mobile":
		return "Mobile", "pnpm dev:llm", nil
	case "desktop":
		return "Desktop", "pnpm dev:lld", nil
	default:
		return "", "", fmt.Errorf("invalid platform selection")
	}
}

func selectParameters(availableParams []Parameter) ([]Parameter, error) {
	// Build options for huh multi-select
	var options []huh.Option[string]
	for _, param := range availableParams {
		options = append(options, huh.NewOption(param.Name, param.Name))
	}

	var selectedNames []string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose parameters:").
				Description("Select the parameters you want to use."). // make dynamic based on selected option
				Options(options...).
				Value(&selectedNames),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		return []Parameter{}, fmt.Errorf("parameter selection cancelled")
	}

	// Convert selected names back to Parameter structs
	var selectedParams []Parameter
	for _, selectedName := range selectedNames {
		for _, param := range availableParams {
			if param.Name == selectedName {
				selectedParams = append(selectedParams, param)
				break
			}
		}
	}

	return selectedParams, nil
}

func inputPresetName(existingPresets []Preset) (string, error) {
	var name string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter preset name:").
				Placeholder("e.g., 'Mobile Dev', 'Desktop Testing'").
				Value(&name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("preset name cannot be empty")
					}
					for _, p := range existingPresets {
						if p.Name == s {
							return fmt.Errorf("preset '%s' already exists", s)
						}
					}
					return nil
				}),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		return "", fmt.Errorf("preset name input cancelled")
	}
	
	return strings.TrimSpace(name), nil
}

// Prefilled versions for editing

func selectPlatformWithDefault(defaultPlatform string) (string, string, error) {
	var options []huh.Option[string]
	options = append(options, huh.NewOption("Mobile", "mobile"))
	options = append(options, huh.NewOption("Desktop", "desktop"))

	var selected string = defaultPlatform // Set default
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Start Ledger-Live for:").
				Options(options...).
				Value(&selected),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		return "", "", fmt.Errorf("platform selection cancelled")
	}

	switch selected {
	case "mobile":
		return "Mobile", "pnpm dev:llm", nil
	case "desktop":
		return "Desktop", "pnpm dev:lld", nil
	default:
		return "", "", fmt.Errorf("invalid platform selection")
	}
}

func selectParametersWithDefault(availableParams []Parameter, selectedParameterNames []string) ([]Parameter, error) {
	// Build options for huh multi-select
	var options []huh.Option[string]
	for _, param := range availableParams {
		option := huh.NewOption(param.Name, param.Name)
		// Check if this parameter should be preselected
		for _, selectedName := range selectedParameterNames {
			if param.Name == selectedName {
				option = option.Selected(true)
				break
			}
		}
		options = append(options, option)
	}

	var selectedNames []string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose parameters:").
				Description("Select additional parameters.").
				Options(options...).
				Value(&selectedNames),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		return []Parameter{}, fmt.Errorf("parameter selection cancelled")
	}

	// Convert selected names back to Parameter structs
	var selectedParams []Parameter
	for _, selectedName := range selectedNames {
		for _, param := range availableParams {
			if param.Name == selectedName {
				selectedParams = append(selectedParams, param)
				break
			}
		}
	}

	return selectedParams, nil
}

func inputPresetNameWithDefault(existingPresets []Preset, currentName string) (string, error) {
	var name string = currentName // Set default
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter preset name:").
				Placeholder("e.g., 'Mobile Dev', 'Desktop Testing'").
				Value(&name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("preset name cannot be empty")
					}
					// Allow keeping the same name
					if s == currentName {
						return nil
					}
					for _, p := range existingPresets {
						if p.Name == s {
							return fmt.Errorf("preset '%s' already exists", s)
						}
					}
					return nil
				}),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		return "", fmt.Errorf("preset name input cancelled")
	}
	
	return strings.TrimSpace(name), nil
}
