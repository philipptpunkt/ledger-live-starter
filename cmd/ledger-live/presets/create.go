package presets

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// CreatePreset creates a new preset from the main flow
func CreatePreset() {
	config, _ := loadConfigWithDefaults()
	
	preset, err := createPresetCore(config)
	if err != nil {
		return
	}

	// Show post-creation menu for main flow
	showPostCreationMenu(*preset, config)
}

// CreatePresetFromManagement creates a new preset from the management menu
func CreatePresetFromManagement(config *setup.Config) {
	fmt.Println(TitleText("Create New Preset"))
	fmt.Println()
	
	preset, err := createPresetCore(config)
	if err != nil {
		return
	}

	// Show management post-creation menu
	showManagementPostCreationMenu(*preset, config)
}

// createPresetCore contains the shared preset creation logic
func createPresetCore(config *setup.Config) (*setup.Preset, error) {
	// Step 1: Get preset name
	presetName, err := InputPresetName(config.Presets)
	if err != nil {
		fmt.Printf("%s %s\n", ErrorText("Error:"), NormalText(err.Error()))
		return nil, err
	}

	// Step 2: Platform selection
	_, platformCommand, err := SelectPlatform()
	if err != nil {
		fmt.Printf("%s %s\n", ErrorText("Error:"), NormalText(err.Error()))
		return nil, err
	}

	// Convert platform name to platform key
	platformKey := convertPlatformToKey(platformCommand)

	// Step 3: Parameter selection
	selectedParams, err := SelectParameters(config.Parameters)
	if err != nil {
		fmt.Printf("%s %s\n", ErrorText("Error:"), NormalText(err.Error()))
		return nil, err
	}

	// Convert selected parameters to parameter names
	parameterNames := extractParameterNames(selectedParams)

	// Step 4: Create new preset
	newPreset := setup.Preset{
		Name:       presetName,
		Platform:   platformKey,
		Parameters: parameterNames,
	}

	// Step 5: Save to config
	config.Presets = append(config.Presets, newPreset)
	
	err = setup.SaveConfig(config)
	if err != nil {
		fmt.Printf("%s %s %s\n", ErrorText("Error:"), NormalText("Error saving preset:"), NormalText(err.Error()))
		return nil, err
	}

	fmt.Printf("%s %s '%s' %s\n\n", SuccessText("Success:"), NormalText("Preset"), HighlightText(presetName), NormalText("created successfully!"))
	
	// Show preset summary
	displayPresetSummary(newPreset)
	
	return &newPreset, nil
}

// showPostCreationMenu handles post-creation options for main flow
func showPostCreationMenu(createdPreset setup.Preset, config *setup.Config) {
	var options []huh.Option[string]
	options = append(options, huh.NewOption("Run this preset now", "run"))
	options = append(options, huh.NewOption("Add another preset", "add"))
	options = append(options, huh.NewOption("Exit", "exit"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do next?").
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
	case "run":
		// Execute the newly created preset
		fmt.Printf("%s %s %s\n", SuccessText("✓"), NormalText("Starting preset:"), HighlightText(createdPreset.Name))
		cmdInfo := BuildPresetCommand(&createdPreset, config)
		ExecuteCommand(cmdInfo)
	case "add":
		// Create another preset
		CreatePreset()
	case "exit":
		fmt.Printf("%s %s\n", SuccessText("✓"), NormalText("Goodbye!"))
		os.Exit(0)
	}
}

// showManagementPostCreationMenu handles post-creation options for management flow
func showManagementPostCreationMenu(createdPreset setup.Preset, config *setup.Config) {
	var options []huh.Option[string]
	options = append(options, huh.NewOption("Run this preset now", "run"))
	options = append(options, huh.NewOption("Add another preset", "add"))
	options = append(options, huh.NewOption("Back to preset management", "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do next?").
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
	case "run":
		// Execute the newly created preset
		fmt.Printf("%s %s %s\n", SuccessText("✓"), NormalText("Starting preset:"), HighlightText(createdPreset.Name))
		cmdInfo := BuildPresetCommand(&createdPreset, config)
		ExecuteCommand(cmdInfo)
	case "add":
		// Create another preset
		CreatePresetFromManagement(config)
	case "back":
		// Reload config and show management menu
		newConfig, err := setup.LoadConfig()
		if err != nil {
			fmt.Printf("%s %s\n", WarningText("Warning:"), NormalText(fmt.Sprintf("Could not reload config (%v)", err)))
			return
		}
		ShowManagementMenu(newConfig)
	}
}
