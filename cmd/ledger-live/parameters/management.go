package parameters

import (
	"github.com/charmbracelet/huh"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// ShowManagementMenu displays the main parameter management menu
func ShowManagementMenu(config *setup.Config) {
	var options []huh.Option[string]
	
	// Always available options
	options = append(options, huh.NewOption("Add new parameter", "add"))
	
	// Only show edit/delete if parameters exist
	if len(config.Parameters) > 0 {
		options = append(options, huh.NewOption("Edit parameters", "edit"))
		options = append(options, huh.NewOption("Delete parameters", "delete"))
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
		AddNewParameter(config)
	case "edit":
		ShowEditParametersMenu(config)
	case "delete":
		ShowDeleteParametersMenu(config)
	case "show":
		ShowAllParameters(config)
	case "back":
		// Load config and navigate back to more menu
		config, err := setup.LoadConfig()
		if err != nil {
			config = setup.GetDefaultConfig()
		}
		ShowMoreMenu(config)
		return
	}
}
