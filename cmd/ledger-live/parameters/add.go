package parameters

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// AddNewParameter handles the creation of a new parameter
func AddNewParameter(config *setup.Config) {
	fmt.Println(TitleText("Add new parameter"))
	fmt.Println()

	// Get parameter name
	name, err := getParameterName("", config.Parameters)
	if err != nil {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Parameter name input cancelled"))
		return
	}

	// Get environment variable
	envVar, err := getEnvironmentVariable("")
	if err != nil {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Environment variable input cancelled"))
		return
	}

	// Get description
	description, err := getParameterDescription("")
	if err != nil {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("Description input cancelled"))
		return
	}

	// Create new parameter
	newParam := setup.Parameter{
		Name:        strings.TrimSpace(name),
		EnvVar:      strings.TrimSpace(envVar),
		Description: strings.TrimSpace(description),
	}

	config.Parameters = append(config.Parameters, newParam)

	// Save config
	err = saveConfigWithError(config)
	if err != nil {
		return
	}

	fmt.Printf("%s %s '%s' %s\n\n", SuccessText("Success:"), NormalText("Parameter"), HighlightText(newParam.Name), NormalText("added successfully!"))
	
	// Show the parameter management menu again
	ShowManagementMenu(config)
}

// getParameterName gets parameter name with validation
func getParameterName(currentValue string, existingParams []setup.Parameter) (string, error) {
	var name string = currentValue
	
	nameForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter parameter name:").
				Placeholder("e.g., 'Skip onboarding', 'Debug mode'").
				Value(&name).
				Validate(func(s string) error {
					return validateParameterName(s, existingParams, currentValue)
				}),
		),
	)

	err := RunStyledForm(nameForm)
	if err != nil {
		return "", err
	}

	return name, nil
}

// getEnvironmentVariable gets environment variable with validation
func getEnvironmentVariable(currentValue string) (string, error) {
	var envVar string = currentValue
	
	envVarForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter environment variable:").
				Placeholder("e.g., 'SKIP_ONBOARDING=1', 'DEBUG_MODE=true'").
				Value(&envVar).
				Validate(validateEnvironmentVariable),
		),
	)

	err := RunStyledForm(envVarForm)
	if err != nil {
		return "", err
	}

	return envVar, nil
}

// getParameterDescription gets parameter description (optional)
func getParameterDescription(currentValue string) (string, error) {
	var description string = currentValue
	
	descForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter description:").
				Placeholder("e.g., 'Skip the onboarding process on mobile'").
				Value(&description),
		),
	)

	err := RunStyledForm(descForm)
	if err != nil {
		return "", err
	}

	return description, nil
}
