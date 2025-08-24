package parameters

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// ShowEditParametersMenu displays the parameter selection menu for editing
func ShowEditParametersMenu(config *setup.Config) {
	if len(config.Parameters) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No parameters found to edit."))
		return
	}

	options := createParameterOptionsFromList(config.Parameters)
	options = append(options, huh.NewOption("Back", "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a parameter to edit:").
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

	EditSingleParameter(selected, config)
}

// EditSingleParameter handles editing of a specific parameter
func EditSingleParameter(paramName string, config *setup.Config) {
	// Find the parameter
	paramIndex, currentParam := findParameterByName(paramName, config.Parameters)
	
	if currentParam == nil {
		fmt.Printf("%s %s '%s' %s\n", ErrorText("Error:"), NormalText("Parameter"), HighlightText(paramName), NormalText("not found"))
		return
	}

	fmt.Printf("%s %s\n\n", TitleText("Editing parameter:"), NormalText(currentParam.Name))

	// Pre-fill with current values
	name := currentParam.Name
	envVar := currentParam.EnvVar
	description := currentParam.Description

	// Create a single form with all fields pre-filled
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Parameter name:").
				Value(&name).
				Validate(func(s string) error {
					return validateParameterName(s, config.Parameters, currentParam.Name)
				}),
			huh.NewInput().
				Title("Environment variable:").
				Value(&envVar).
				Validate(validateEnvironmentVariable),
			huh.NewInput().
				Title("Description:").
				Value(&description),
		),
	)

	err := RunStyledForm(form)
	if err != nil {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Parameter editing cancelled"))
		return
	}

	// Update the parameter with new values
	config.Parameters[paramIndex].Name = strings.TrimSpace(name)
	config.Parameters[paramIndex].EnvVar = strings.TrimSpace(envVar)
	config.Parameters[paramIndex].Description = strings.TrimSpace(description)

	// Save config
	err = saveConfigWithError(config)
	if err != nil {
		return
	}

	fmt.Printf("%s %s '%s' %s\n\n", SuccessText("Success:"), NormalText("Parameter"), HighlightText(strings.TrimSpace(name)), NormalText("updated successfully!"))
	
	// Return to parameter management menu
	ShowManagementMenu(config)
}
