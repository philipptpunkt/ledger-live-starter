package parameters

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// ShowDeleteParametersMenu displays the parameter selection menu for deletion
func ShowDeleteParametersMenu(config *setup.Config) {
	if len(config.Parameters) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No parameters available to delete."))
		return
	}

	options := createParameterOptionsFromList(config.Parameters)

	var selectedParameters []string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Choose parameters:").
				DescriptionFunc(func() string {
					if len(selectedParameters) == 0 {
						return "Select parameter(s) to delete (use no selection and ENTER to exit)"
					}
					if len(selectedParameters) == 1 {
						// Show description of the single selected parameter
						for _, param := range config.Parameters {
							if param.Name == selectedParameters[0] {
								if param.Description != "" {
									return param.Description
								}
								return "- no description available -"
							}
						}
					}
					// Multiple parameters selected
					return fmt.Sprintf("%d parameters selected for deletion", len(selectedParameters))
				}, &selectedParameters).
				Options(options...).
				Value(&selectedParameters),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		fmt.Printf("\n%s %s\n", ErrorText("❌"), NormalText("Selection cancelled"))
		ShowManagementMenu(config)
		return
	}

	if len(selectedParameters) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No parameters selected for deletion."))
		ShowManagementMenu(config)
		return
	}

	// Show confirmation dialog
	confirmParameterDeletion(selectedParameters, config)
}

// confirmParameterDeletion shows confirmation dialog for parameter deletion
func confirmParameterDeletion(paramNames []string, config *setup.Config) {
	var confirm bool
	
	// Build confirmation message
	message := fmt.Sprintf("Are you sure you want to delete %d parameter(s)?", len(paramNames))
	if len(paramNames) == 1 {
		message = fmt.Sprintf("Are you sure you want to delete '%s'?", paramNames[0])
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
		// Delete the selected parameters
		deleteMultipleParameters(paramNames, config)
	} else {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Deletion cancelled."))
		ShowManagementMenu(config)
	}
}

// deleteMultipleParameters removes multiple parameters from config
func deleteMultipleParameters(paramNames []string, config *setup.Config) {
	// Create a set of parameter names for faster lookup
	toDelete := make(map[string]bool)
	for _, name := range paramNames {
		toDelete[name] = true
	}

	// Filter out the parameters to delete
	var remainingParameters []setup.Parameter
	var deletedCount int
	for _, param := range config.Parameters {
		if toDelete[param.Name] {
			deletedCount++
		} else {
			remainingParameters = append(remainingParameters, param)
		}
	}

	// Update config
	config.Parameters = remainingParameters

	// Save config
	err := saveConfigWithError(config)
	if err != nil {
		return
	}

	// Show success message
	if deletedCount == 1 {
		fmt.Printf("%s %s '%s' %s\n", SuccessText("Success:"), NormalText("Parameter"), HighlightText(paramNames[0]), NormalText("deleted successfully."))
	} else {
		fmt.Printf("%s %s %s\n", SuccessText("Success:"), HighlightText(fmt.Sprintf("%d", deletedCount)), NormalText("parameters deleted successfully."))
	}

	// Return to management menu
	ShowManagementMenu(config)
}
