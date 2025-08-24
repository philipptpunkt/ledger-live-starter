package parameters

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// ShowAllParameters displays all parameters with action options
func ShowAllParameters(config *setup.Config) {
	if len(config.Parameters) == 0 {
		fmt.Printf("%s %s\n", InfoTextTitle("Info:"), NormalText("No parameters defined yet."))
		fmt.Println()
	} else {
		displayParametersList(config.Parameters)
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
		ShowEditParametersMenu(config)
	case "back":
		ShowManagementMenu(config)
	}
}

// displayParametersList shows formatted list of all parameters
func displayParametersList(parameters []setup.Parameter) {
	for i, param := range parameters {
		fmt.Printf("%s %s %d:\n", TitleText("•"), NormalText("Parameter"), i+1)
		fmt.Printf("   %s %s\n", InfoTextTitle("Name:"), HighlightText(param.Name))
		fmt.Printf("   %s %s\n", InfoTextTitle("Environment Variable:"), HighlightText(param.EnvVar))
		if param.Description != "" {
			fmt.Printf("   %s %s\n", InfoTextTitle("Description:"), NormalText(param.Description))
		}
		fmt.Println()
	}
}
