package presets

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// ShowDeletePresetsMenu displays the preset selection menu for deletion
func ShowDeletePresetsMenu(config *setup.Config) {
	if len(config.Presets) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No presets available to delete."))
		return
	}

	var options []huh.Option[string]
	for _, preset := range config.Presets {
		options = append(options, huh.NewOption(preset.Name, preset.Name))
	}

	var selectedPresets []string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose presets:").
				Description("Select preset(s) to delete (use no selection and ENTER to exit):").
				Options(options...).
				Value(&selectedPresets),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		ShowCancellationMessage()
		ShowManagementMenu(config)
		return
	}

	if len(selectedPresets) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No presets selected for deletion."))
		ShowManagementMenu(config)
		return
	}

	// Show confirmation dialog
	confirmDeletion(selectedPresets, config)
}

// confirmDeletion shows confirmation dialog for preset deletion
func confirmDeletion(presetNames []string, config *setup.Config) {
	var confirm bool
	
	// Build confirmation message
	message := fmt.Sprintf("Are you sure you want to delete %d preset(s)?", len(presetNames))
	if len(presetNames) == 1 {
		message = fmt.Sprintf("Are you sure you want to delete '%s'?", presetNames[0])
	}
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(message).
				Affirmative("Delete").
				Negative("Cancel").
				Value(&confirm),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		ShowConfirmationCancelledMessage()
		ShowManagementMenu(config)
		return
	}

	if confirm {
		// Delete the selected presets
		deleteMultiplePresets(presetNames, config)
	} else {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Deletion cancelled."))
		ShowManagementMenu(config)
	}
}

// deleteMultiplePresets removes multiple presets from config
func deleteMultiplePresets(presetNames []string, config *setup.Config) {
	// Create a set of preset names for faster lookup
	toDelete := make(map[string]bool)
	for _, name := range presetNames {
		toDelete[name] = true
	}

	// Filter out the presets to delete
	var remainingPresets []setup.Preset
	var deletedCount int
	for _, preset := range config.Presets {
		if toDelete[preset.Name] {
			deletedCount++
		} else {
			remainingPresets = append(remainingPresets, preset)
		}
	}

	// Update config
	config.Presets = remainingPresets

	// Save config
	err := saveConfigWithError(config)
	if err != nil {
		return
	}

	// Show success message
	if deletedCount == 1 {
		fmt.Printf("%s %s '%s' %s\n", SuccessText("Success:"), NormalText("Preset"), HighlightText(presetNames[0]), NormalText("deleted successfully."))
	} else {
		fmt.Printf("%s %s %s %s\n", SuccessText("Success:"), HighlightText(fmt.Sprintf("%d", deletedCount)), NormalText("presets deleted successfully."), NormalText(""))
	}

	// Return to management menu
	ShowManagementMenu(config)
}
