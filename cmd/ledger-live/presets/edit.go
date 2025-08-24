package presets

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// EditPresets entry point for preset editing
func EditPresets() {
	// Load configuration
	config, err := setup.LoadConfig()
	if err != nil {
		fmt.Printf("%s %s\n", WarningText("Warning:"), NormalText(fmt.Sprintf("Could not load config.json (%v)", err)))
		return
	}

	if len(config.Presets) == 0 {
		// No presets exist, offer to create one
		showNoPresetsEditMenu()
		return
	}

	// Show preset management menu
	ShowManagementMenu(config)
}

// showNoPresetsEditMenu displays menu when no presets exist for editing
func showNoPresetsEditMenu() {
	fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No presets found to edit."))
	fmt.Println()

	var options []huh.Option[string]
	options = append(options, huh.NewOption("Create new preset", "create"))
	options = append(options, huh.NewOption("Back", "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
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
		CreatePreset()
	case "back":
		return
	}
}

// ShowEditPresetsMenu displays the preset selection menu for editing
func ShowEditPresetsMenu(config *setup.Config) {
	if len(config.Presets) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No presets available to edit."))
		return
	}

	var options []huh.Option[string]
	for _, preset := range config.Presets {
		options = append(options, huh.NewOption(preset.Name, preset.Name))
	}
	options = append(options, huh.NewOption("Back", "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a preset to edit:").
				Options(options...).
				Value(&selected),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		ShowCancellationMessage()
		return
	}

	if selected == "back" {
		ShowManagementMenu(config)
		return
	}

	// Edit the selected preset
	EditSinglePreset(selected, config)
}

// EditSinglePreset handles editing of a specific preset
func EditSinglePreset(presetName string, config *setup.Config) {
	// Find the preset
	presetIndex, currentPreset := findPresetByName(presetName, config.Presets)
	
	if currentPreset == nil {
		fmt.Printf("%s %s '%s' %s\n", ErrorText("❌"), NormalText("Preset"), HighlightText(presetName), NormalText("not found"))
		ShowEditPresetsMenu(config)
		return
	}

	fmt.Printf("%s %s %s\n\n", TitleText("✏️"), NormalText("Editing preset:"), HighlightText(currentPreset.Name))

	// Create a single form with all fields prefilled
	var newName string = currentPreset.Name
	var newPlatform string = currentPreset.Platform
	var selectedParameterNames []string = currentPreset.Parameters

	// Create platform options with current selection
	var platformOptions []huh.Option[string]
	platformOptions = append(platformOptions, huh.NewOption("Mobile", "mobile"))
	platformOptions = append(platformOptions, huh.NewOption("Desktop", "desktop"))

	// Create parameter options with current selections
	var parameterOptions []huh.Option[string]
	for _, param := range config.Parameters {
		option := huh.NewOption(param.Name, param.Name)
		// Check if this parameter should be preselected
		for _, selectedName := range currentPreset.Parameters {
			if param.Name == selectedName {
				option = option.Selected(true)
				break
			}
		}
		parameterOptions = append(parameterOptions, option)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Preset name:").
				Value(&newName).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("preset name cannot be empty")
					}
					// Allow keeping the same name
					if s == currentPreset.Name {
						return nil
					}
					for _, p := range config.Presets {
						if p.Name == s {
							return fmt.Errorf("preset '%s' already exists", s)
						}
					}
					return nil
				}),
			
			huh.NewSelect[string]().
				Title("Platform:").
				Options(platformOptions...).
				Value(&newPlatform),
				
			huh.NewMultiSelect[string]().
				Title("Parameters (use SPACE to toggle):").
				Options(parameterOptions...).
				Value(&selectedParameterNames),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		fmt.Printf("\n%s %s\n", ErrorText("❌"), NormalText("Edit cancelled"))
		ShowEditPresetsMenu(config)
		return
	}

	// Update the preset with new values
	config.Presets[presetIndex].Name = strings.TrimSpace(newName)
	config.Presets[presetIndex].Platform = newPlatform
	config.Presets[presetIndex].Parameters = selectedParameterNames

	// Save changes
	err = saveConfigWithError(config)
	if err != nil {
		ShowEditPresetsMenu(config)
		return
	}

	fmt.Printf("%s %s '%s' %s\n", SuccessText("Success:"), NormalText("Preset"), HighlightText(newName), NormalText("updated successfully!"))
	
	// Return to edit presets menu
	ShowEditPresetsMenu(config)
}
