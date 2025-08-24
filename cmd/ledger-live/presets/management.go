package presets

import (
	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// ShowManagementMenu displays the main preset management menu
func ShowManagementMenu(config *setup.Config) {
	var options []huh.Option[string]
	
	// Add new preset option first
	options = append(options, huh.NewOption("Add new preset", "add"))
	
	// Add edit presets option
	options = append(options, huh.NewOption("Edit presets", "edit"))
	
	// Add delete presets option
	options = append(options, huh.NewOption("Delete presets", "delete"))
	
	// Add back option last
	options = append(options, huh.NewOption("Back", "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a preset to manage:").
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
	if selected == "back" {
		// Navigate back to more menu
		ShowMoreMenu(config)
		return
	} else if selected == "add" {
		// Create new preset from management menu
		CreatePresetFromManagement(config)
	} else if selected == "edit" {
		// Show edit presets menu
		ShowEditPresetsMenu(config)
	} else if selected == "delete" {
		// Show delete presets menu
		ShowDeletePresetsMenu(config)
	}
}
