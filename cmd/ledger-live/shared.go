package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// Shared UI functions that can be reused across different flows

func selectPlatform() (string, string, error) {
	var options []huh.Option[string]
	options = append(options, huh.NewOption(BoldColorText("Mobile", Blue), "mobile"))
	options = append(options, huh.NewOption(BoldColorText("Desktop", Purple), "desktop"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Start Ledger-Live for:").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
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
				Title("Select additional parameters (use SPACE to toggle, ENTER to continue):").
				Options(options...).
				Value(&selectedNames),
		),
	)

	err := form.Run()
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

	err := form.Run()
	if err != nil {
		return "", fmt.Errorf("preset name input cancelled")
	}
	
	return strings.TrimSpace(name), nil
}
