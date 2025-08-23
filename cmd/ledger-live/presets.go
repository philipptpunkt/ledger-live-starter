package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

func createPreset() {
	fmt.Println(BoldColorText("Create New Preset", Green))
	fmt.Println()

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("%s Warning: Could not load config.json (%v), using defaults\n\n", BoldColorText("!", Yellow), err)
		config = getDefaultConfig()
	}

	// Step 1: Get preset name
	presetName, err := inputPresetName(config.Presets)
	if err != nil {
		fmt.Printf("%s %v\n", ColorText("Error:", Red), err)
		return
	}

	// Step 2: Platform selection
	_, platformCommand, err := selectPlatform()
	if err != nil {
		fmt.Printf("%s %v\n", ColorText("Error:", Red), err)
		return
	}

	// Convert platform name to platform key
	var platformKey string
	switch platformCommand {
	case "pnpm dev:llm":
		platformKey = "mobile"
	case "pnpm dev:lld":
		platformKey = "desktop"
	default:
		platformKey = "mobile"
	}

	// Step 3: Parameter selection
	selectedParams, err := selectParameters(config.Parameters)
	if err != nil {
		fmt.Printf("%s %v\n", ColorText("Error:", Red), err)
		return
	}

	// Convert selected parameters to parameter names
	var parameterNames []string
	for _, param := range selectedParams {
		parameterNames = append(parameterNames, param.Name)
	}

	// Step 4: Create new preset
	newPreset := Preset{
		Name:       presetName,
		Platform:   platformKey,
		Parameters: parameterNames,
	}

	// Step 5: Save to config
	config.Presets = append(config.Presets, newPreset)
	
	err = saveConfig(config)
	if err != nil {
		fmt.Printf("%s Error saving preset: %v\n", ColorText("Error:", Red), err)
		return
	}

	fmt.Printf("%s Preset '%s' created successfully!\n\n", ColorText("Success:", Green), presetName)
	
	// Show preset summary
	fmt.Println("📋 Preset Summary:")
	fmt.Printf("   Name: %s\n", newPreset.Name)
	fmt.Printf("   Platform: %s\n", strings.Title(newPreset.Platform))
	if len(newPreset.Parameters) > 0 {
		fmt.Printf("   Parameters: %s\n", strings.Join(newPreset.Parameters, ", "))
	} else {
		fmt.Println("   Parameters: None")
	}
	fmt.Println()

	// Step 6: Post-creation menu
	showPostCreationMenu(newPreset, config)
}

func showPostCreationMenu(createdPreset Preset, config *Config) {
	var options []huh.Option[string]
	options = append(options, huh.NewOption(BoldColorText("Run this preset now", Green), "run"))
	options = append(options, huh.NewOption(BoldColorText("Add another preset", Blue), "add"))
	options = append(options, huh.NewOption(ColorText("Exit", Red), "exit"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do next?").
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
	case "run":
		// Execute the newly created preset
		fmt.Printf("✅ Starting preset: %s\n", createdPreset.Name)
		cmdInfo := buildPresetCommand(&createdPreset, config)
		executeCommand(cmdInfo)
	case "add":
		// Create another preset
		createPreset()
	case "exit":
		fmt.Printf("%s Goodbye!\n", BoldColorText("✓", Green))
		os.Exit(0)
	}
}

func editPresets() {
	fmt.Println(BoldColorText("Edit Presets", Blue))
	fmt.Println()

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("%s Warning: Could not load config.json (%v)\n", BoldColorText("!", Yellow), err)
		return
	}

	if len(config.Presets) == 0 {
		// No presets exist, offer to create one
		showNoPresetsEditMenu()
		return
	}

	// Show preset management menu
	showPresetManagementMenu(config)
}

func showNoPresetsEditMenu() {
	fmt.Printf("%s No presets found to edit.\n", ColorText("Info:", Yellow))
	fmt.Println()

	var options []huh.Option[string]
	options = append(options, huh.NewOption(BoldColorText("Create new preset", Green), "create"))
	options = append(options, huh.NewOption(ColorText("Back", Red), "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to do?").
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
	case "back":
		return
	}
}

func showPresetManagementMenu(config *Config) {
	var options []huh.Option[string]
	
	// Add each preset as an edit option
	for _, preset := range config.Presets {
		options = append(options, huh.NewOption(fmt.Sprintf("%s Edit '%s'", BoldColorText("Edit", Blue), preset.Name), "edit:"+preset.Name))
	}
	
	// Add delete options
	for _, preset := range config.Presets {
		options = append(options, huh.NewOption(fmt.Sprintf("%s Delete '%s'", BoldColorText("Delete", Red), preset.Name), "delete:"+preset.Name))
	}
	
	// Add back option
	options = append(options, huh.NewOption(ColorText("Back to main menu", Red), "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a preset to manage:").
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
	if selected == "back" {
		return
	} else if strings.HasPrefix(selected, "edit:") {
		presetName := strings.TrimPrefix(selected, "edit:")
		editSinglePreset(presetName, config)
	} else if strings.HasPrefix(selected, "delete:") {
		presetName := strings.TrimPrefix(selected, "delete:")
		deletePreset(presetName, config)
	}
}

func editSinglePreset(presetName string, config *Config) {
	// Find the preset
	var presetIndex = -1
	var currentPreset *Preset
	
	for i, preset := range config.Presets {
		if preset.Name == presetName {
			presetIndex = i
			currentPreset = &preset
			break
		}
	}

	if currentPreset == nil {
		fmt.Printf("❌ Preset '%s' not found\n", presetName)
		return
	}

	fmt.Printf("✏️  Editing preset: %s\n\n", currentPreset.Name)

	// Show what to edit
	var options []huh.Option[string]
	options = append(options, huh.NewOption("📝 Edit name", "name"))
	options = append(options, huh.NewOption("📱 Change platform", "platform"))
	options = append(options, huh.NewOption("⚙️  Edit parameters", "parameters"))
	options = append(options, huh.NewOption("← Back", "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("What would you like to edit?").
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
	case "name":
		newName, err := inputPresetName(config.Presets)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
		if newName != "" && newName != currentPreset.Name {
			config.Presets[presetIndex].Name = newName
			saveAndConfirm(config, fmt.Sprintf("✅ Preset name changed to '%s'", newName))
		}
	case "platform":
		platformName, platformCommand, err := selectPlatform()
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
		if platformName != "" {
			var platformKey string
			switch platformCommand {
			case "pnpm dev:llm":
				platformKey = "mobile"
			case "pnpm dev:lld":
				platformKey = "desktop"
			default:
				platformKey = "mobile"
			}
			config.Presets[presetIndex].Platform = platformKey
			saveAndConfirm(config, fmt.Sprintf("✅ Platform changed to %s", strings.Title(platformKey)))
		}
	case "parameters":
		selectedParams, err := selectParameters(config.Parameters)
		if err != nil {
			fmt.Printf("❌ %v\n", err)
			return
		}
		var parameterNames []string
		for _, param := range selectedParams {
			parameterNames = append(parameterNames, param.Name)
		}
		config.Presets[presetIndex].Parameters = parameterNames
		saveAndConfirm(config, "✅ Parameters updated")
	case "back":
		return
	}
	
	// After editing, show the menu again
	showPresetManagementMenu(config)
}

func deletePreset(presetName string, config *Config) {
	// Find the preset index
	var presetIndex = -1
	
	for i, preset := range config.Presets {
		if preset.Name == presetName {
			presetIndex = i
			break
		}
	}

	if presetIndex == -1 {
		fmt.Printf("❌ Preset '%s' not found\n", presetName)
		return
	}

	// Confirm deletion
	var confirm bool
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Are you sure you want to delete '%s'?", presetName)).
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("\n❌ Deletion cancelled")
		return
	}

	if confirm {
		// Remove the preset from the slice
		config.Presets = append(config.Presets[:presetIndex], config.Presets[presetIndex+1:]...)
		saveAndConfirm(config, fmt.Sprintf("✅ Preset '%s' deleted", presetName))
		
		// If no presets left, go back to main menu
		if len(config.Presets) == 0 {
			fmt.Println("No presets remaining. Returning to main menu.")
			return
		}
	}
	
	// Show the menu again
	showPresetManagementMenu(config)
}

func saveAndConfirm(config *Config, message string) {
	err := saveConfig(config)
	if err != nil {
		fmt.Printf("❌ Error saving changes: %v\n", err)
	} else {
		fmt.Println(message)
	}
}

func saveConfig(config *Config) error {
	return saveConfigToPath(config, getConfigPath())
}
