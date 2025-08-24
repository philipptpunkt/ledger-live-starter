package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// Parameter management functionality

func editParameters() {
	// Load configuration
	config, err := setup.LoadConfig()
	if err != nil {
		fmt.Printf("%s %s\n\n", WarningText("Warning:"), NormalText(fmt.Sprintf("Could not load config.json (%v), using defaults", err)))
		config = setup.GetDefaultConfig()
	}

	// Show parameter management menu
	showParameterManagementMenu(config)
}

func showParameterManagementMenu(config *Config) {
	var options []huh.Option[string]
	
	// Always available options
	options = append(options, huh.NewOption("Add new parameter", "add"))
	
	// Only show edit/delete if parameters exist
	if len(config.Parameters) > 0 {
		options = append(options, huh.NewOption("Edit parameter", "edit"))
		options = append(options, huh.NewOption("Delete parameter", "delete"))
	}
	
	options = append(options, huh.NewOption("Show all parameters", "show"))
	options = append(options, huh.NewOption("Back", "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose an option:").
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
	case "add":
		addNewParameter(config)
	case "edit":
		editSingleParameterMenu(config)
	case "delete":
		showDeleteParametersMenu(config)
	case "show":
		showAllParameters(config)
	case "back":
		// Load config and navigate back to more menu
		config, err := setup.LoadConfig()
		if err != nil {
			config = setup.GetDefaultConfig()
		}
		showMoreMenu(config)
		return
	}
}

func addNewParameter(config *Config) {
	fmt.Println(TitleText("Add New Parameter"))
	fmt.Println()

	// Get parameter name
	var name string
	nameForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter parameter name:").
				Placeholder("e.g., 'Skip onboarding', 'Debug mode'").
				Value(&name).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("parameter name cannot be empty")
					}
					for _, p := range config.Parameters {
						if p.Name == s {
							return fmt.Errorf("parameter '%s' already exists", s)
						}
					}
					return nil
				}),
		),
	)

	err := RunStyledForm(nameForm)
	if err != nil {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Parameter name input cancelled"))
		return
	}

	// Get environment variable
	var envVar string
	envVarForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter environment variable:").
				Placeholder("e.g., 'SKIP_ONBOARDING=1', 'DEBUG_MODE=true'").
				Value(&envVar).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("environment variable cannot be empty")
					}
					if !strings.Contains(s, "=") {
						return fmt.Errorf("environment variable must include '=' (e.g., VAR_NAME=value)")
					}
					return nil
				}),
		),
	)

	err = RunStyledForm(envVarForm)
	if err != nil {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Environment variable input cancelled"))
		return
	}

	// Get description
	var description string
	descForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter description:").
				Placeholder("e.g., 'Skip the onboarding process on mobile'").
				Value(&description),
		),
	)

	err = RunStyledForm(descForm)
	if err != nil {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Description input cancelled"))
		return
	}

	// Add new parameter
	newParam := Parameter{
		Name:        strings.TrimSpace(name),
		EnvVar:      strings.TrimSpace(envVar),
		Description: strings.TrimSpace(description),
	}

	config.Parameters = append(config.Parameters, newParam)

	// Save config
	err = setup.SaveConfig(config)
	if err != nil {
		fmt.Printf("%s %s %s\n", ErrorText("Error:"), NormalText("Error saving parameter:"), NormalText(err.Error()))
		return
	}

	fmt.Printf("%s %s '%s' %s\n\n", SuccessText("Success:"), NormalText("Parameter"), HighlightText(newParam.Name), NormalText("added successfully!"))
	
	// Show the parameter management menu again
	showParameterManagementMenu(config)
}

func editSingleParameterMenu(config *Config) {
	if len(config.Parameters) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No parameters found to edit."))
		return
	}

	var options []huh.Option[string]
	for _, param := range config.Parameters {
		options = append(options, huh.NewOption(param.Name, param.Name))
	}
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
		showParameterManagementMenu(config)
		return
	}

	editSingleParameter(selected, config)
}

func editSingleParameter(paramName string, config *Config) {
	// Find the parameter
	var paramIndex = -1
	var currentParam *Parameter
	
	for i, param := range config.Parameters {
		if param.Name == paramName {
			paramIndex = i
			currentParam = &param
			break
		}
	}

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
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("parameter name cannot be empty")
					}
					for _, p := range config.Parameters {
						if p.Name == s && p.Name != currentParam.Name {
							return fmt.Errorf("parameter '%s' already exists", s)
						}
					}
					return nil
				}),
			huh.NewInput().
				Title("Environment variable:").
				Value(&envVar).
				Validate(func(s string) error {
					if strings.TrimSpace(s) == "" {
						return fmt.Errorf("environment variable cannot be empty")
					}
					if !strings.Contains(s, "=") {
						return fmt.Errorf("environment variable must include '=' (e.g., VAR_NAME=value)")
					}
					return nil
				}),
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
	err = setup.SaveConfig(config)
	if err != nil {
		fmt.Printf("%s %s %s\n", ErrorText("Error:"), NormalText("Error saving parameter:"), NormalText(err.Error()))
		return
	}

	fmt.Printf("%s %s '%s' %s\n\n", SuccessText("Success:"), NormalText("Parameter"), HighlightText(strings.TrimSpace(name)), NormalText("updated successfully!"))
	
	// Return to parameter management menu
	showParameterManagementMenu(config)
}

func showDeleteParametersMenu(config *Config) {
	if len(config.Parameters) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No parameters available to delete."))
		return
	}

	var options []huh.Option[string]
	for _, param := range config.Parameters {
		options = append(options, huh.NewOption(param.Name, param.Name))
	}

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
		showParameterManagementMenu(config)
		return
	}

	if len(selectedParameters) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No parameters selected for deletion."))
		showParameterManagementMenu(config)
		return
	}

	// Show confirmation dialog
	confirmParameterDeletion(selectedParameters, config)
}

func confirmParameterDeletion(paramNames []string, config *Config) {
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
		showParameterManagementMenu(config)
		return
	}

	if confirm {
		// Delete the selected parameters
		deleteMultipleParameters(paramNames, config)
	} else {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Deletion cancelled."))
		showParameterManagementMenu(config)
	}
}

func deleteMultipleParameters(paramNames []string, config *Config) {
	// Create a set of parameter names for faster lookup
	toDelete := make(map[string]bool)
	for _, name := range paramNames {
		toDelete[name] = true
	}

	// Filter out the parameters to delete
	var remainingParameters []Parameter
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
	err := setup.SaveConfig(config)
	if err != nil {
		fmt.Printf("%s %s %s\n", ErrorText("Error:"), NormalText("Error saving changes:"), NormalText(err.Error()))
		return
	}

	// Show success message
	if deletedCount == 1 {
		fmt.Printf("%s %s '%s' %s\n", SuccessText("Success:"), NormalText("Parameter"), HighlightText(paramNames[0]), NormalText("deleted successfully."))
	} else {
		fmt.Printf("%s %s %s\n", SuccessText("Success:"), HighlightText(fmt.Sprintf("%d", deletedCount)), NormalText("parameters deleted successfully."))
	}

	// Return to management menu
	showParameterManagementMenu(config)
}

func showAllParameters(config *Config) {
	if len(config.Parameters) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No parameters defined yet."))
		fmt.Println()
	} else {
		for i, param := range config.Parameters {
			fmt.Printf("%s %s %d:\n", TitleText("•"), NormalText("Parameter"), i+1)
			fmt.Printf("   %s %s\n", InfoTextTitle("Name:"), HighlightText(param.Name))
			fmt.Printf("   %s %s\n", InfoTextTitle("Environment Variable:"), HighlightText(param.EnvVar))
			if param.Description != "" {
				fmt.Printf("   %s %s\n", InfoTextTitle("Description:"), NormalText(param.Description))
			}
			fmt.Println()
		}
	}

	// Show action options
	var options []huh.Option[string]
	
	// Only show Edit option if parameters exist
	if len(config.Parameters) > 0 {
		options = append(options, huh.NewOption("Edit", "edit"))
	}
	
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
	case "edit":
		editSingleParameterMenu(config)
	case "back":
		showParameterManagementMenu(config)
	}
}
