package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// Parameter management functionality

func editParameters() {
	fmt.Println(BoldColorText("Edit Parameters", Yellow))
	fmt.Println()

	// Load configuration
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("%s Warning: Could not load config.json (%v), using defaults\n\n", BoldColorText("!", Yellow), err)
		config = getDefaultConfig()
	}

	// Show parameter management menu
	showParameterManagementMenu(config)
}

func showParameterManagementMenu(config *Config) {
	var options []huh.Option[string]
	
	// Always available options
	options = append(options, huh.NewOption(BoldColorText("Add new parameter", Green), "add"))
	
	// Only show edit/delete if parameters exist
	if len(config.Parameters) > 0 {
		options = append(options, huh.NewOption(BoldColorText("Edit parameter", Blue), "edit"))
		options = append(options, huh.NewOption(BoldColorText("Delete parameter", Red), "delete"))
	}
	
	options = append(options, huh.NewOption(ColorText("Show all parameters", Cyan), "show"))
	options = append(options, huh.NewOption(ColorText("Back", Red), "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose an option:").
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
	case "add":
		addNewParameter(config)
	case "edit":
		editSingleParameterMenu(config)
	case "delete":
		deleteSingleParameterMenu(config)
	case "show":
		showAllParameters(config)
	case "back":
		return
	}
}

func addNewParameter(config *Config) {
	fmt.Println(BoldColorText("Add New Parameter", Green))
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

	err := nameForm.Run()
	if err != nil {
		fmt.Printf("%s Parameter name input cancelled\n", ColorText("Info:", Yellow))
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

	err = envVarForm.Run()
	if err != nil {
		fmt.Printf("%s Environment variable input cancelled\n", ColorText("Info:", Yellow))
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

	err = descForm.Run()
	if err != nil {
		fmt.Printf("%s Description input cancelled\n", ColorText("Info:", Yellow))
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
	err = saveConfig(config)
	if err != nil {
		fmt.Printf("%s Error saving parameter: %v\n", ColorText("Error:", Red), err)
		return
	}

	fmt.Printf("%s Parameter '%s' added successfully!\n\n", ColorText("Success:", Green), newParam.Name)
	
	// Show the parameter management menu again
	showParameterManagementMenu(config)
}

func editSingleParameterMenu(config *Config) {
	if len(config.Parameters) == 0 {
		fmt.Printf("%s No parameters found to edit.\n", ColorText("Info:", Yellow))
		return
	}

	var options []huh.Option[string]
	for _, param := range config.Parameters {
		options = append(options, huh.NewOption(fmt.Sprintf("%s %s", BoldText("Edit"), param.Name), param.Name))
	}
	options = append(options, huh.NewOption(ColorText("Back", Red), "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a parameter to edit:").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("\n❌ Selection cancelled")
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
		fmt.Printf("%s Parameter '%s' not found\n", ColorText("Error:", Red), paramName)
		return
	}

	fmt.Printf("%s %s\n\n", BoldColorText("Editing parameter:", Blue), currentParam.Name)

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

	err := form.Run()
	if err != nil {
		fmt.Printf("%s Parameter editing cancelled\n", ColorText("Info:", Yellow))
		return
	}

	// Update the parameter with new values
	config.Parameters[paramIndex].Name = strings.TrimSpace(name)
	config.Parameters[paramIndex].EnvVar = strings.TrimSpace(envVar)
	config.Parameters[paramIndex].Description = strings.TrimSpace(description)

	// Save config
	err = saveConfig(config)
	if err != nil {
		fmt.Printf("%s Error saving parameter: %v\n", ColorText("Error:", Red), err)
		return
	}

	fmt.Printf("%s Parameter '%s' updated successfully!\n\n", ColorText("Success:", Green), strings.TrimSpace(name))
	
	// Return to parameter management menu
	showParameterManagementMenu(config)
}

func deleteSingleParameterMenu(config *Config) {
	if len(config.Parameters) == 0 {
		fmt.Printf("%s No parameters found to delete.\n", ColorText("Info:", Yellow))
		return
	}

	var options []huh.Option[string]
	for _, param := range config.Parameters {
		options = append(options, huh.NewOption(fmt.Sprintf("%s %s", BoldColorText("Delete", Red), param.Name), param.Name))
	}
	options = append(options, huh.NewOption(ColorText("Back", Red), "back"))

	var selected string
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a parameter to delete:").
				Options(options...).
				Value(&selected),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("\n❌ Selection cancelled")
		return
	}

	if selected == "back" {
		showParameterManagementMenu(config)
		return
	}

	deleteParameter(selected, config)
}

func deleteParameter(paramName string, config *Config) {
	// Find the parameter index
	var paramIndex = -1
	
	for i, param := range config.Parameters {
		if param.Name == paramName {
			paramIndex = i
			break
		}
	}

	if paramIndex == -1 {
		fmt.Printf("%s Parameter '%s' not found\n", ColorText("Error:", Red), paramName)
		return
	}

	// Confirm deletion
	var confirm bool
	
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Are you sure you want to delete '%s'?", paramName)).
				Value(&confirm),
		),
	)

	err := form.Run()
	if err != nil {
		fmt.Println("\n❌ Deletion cancelled")
		return
	}

	if confirm {
		// Remove the parameter from the slice
		config.Parameters = append(config.Parameters[:paramIndex], config.Parameters[paramIndex+1:]...)
		saveAndConfirm(config, fmt.Sprintf("✅ Parameter '%s' deleted", paramName))
	}
	
	// Show the parameter management menu again
	showParameterManagementMenu(config)
}

func showAllParameters(config *Config) {
	fmt.Println(BoldColorText("All Parameters", Cyan))
	fmt.Println()

	if len(config.Parameters) == 0 {
		fmt.Printf("%s No parameters defined yet.\n", ColorText("Info:", Yellow))
		fmt.Println()
	} else {
		for i, param := range config.Parameters {
			fmt.Printf("%s Parameter %d:\n", BoldText("•"), i+1)
			fmt.Printf("   %s %s\n", BoldText("Name:"), param.Name)
			fmt.Printf("   %s %s\n", BoldText("Environment Variable:"), param.EnvVar)
			if param.Description != "" {
				fmt.Printf("   %s %s\n", BoldText("Description:"), param.Description)
			}
			fmt.Println()
		}
	}

	// Show action options
	var options []huh.Option[string]
	
	// Only show Edit option if parameters exist
	if len(config.Parameters) > 0 {
		options = append(options, huh.NewOption(BoldColorText("Edit", Blue), "edit"))
	}
	
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
	case "edit":
		editSingleParameterMenu(config)
	case "back":
		showParameterManagementMenu(config)
	}
}
