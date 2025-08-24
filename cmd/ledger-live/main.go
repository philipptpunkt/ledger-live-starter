package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"ledger-live-starter/cmd/ledger-live/parameters"
	"ledger-live-starter/cmd/ledger-live/presets"
	"ledger-live-starter/cmd/ledger-live/setup"

)

var configPath string // Global config path variable

var rootCmd = &cobra.Command{
	Use:   "ledger-live",
	Short: "A CLI tool to help users start the ledger-live app",
	Long: `Ledger Live Starter is a CLI tool that provides an easy way to start 
the ledger-live application with different configurations and platforms.

You can use it to quickly start mobile or desktop versions with predefined settings.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	// Add global config flag
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config file path (default: ~/.ledger-live/config.json)")
	
	// Add version flag
	var showVersion bool
	rootCmd.Flags().BoolVarP(&showVersion, "version", "V", false, "show version information")
	
	// Override the Run function to handle --version flag
	originalRun := rootCmd.Run
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if showVersion {
			// Just show the version number
			fmt.Printf(Version)
			return
		}
		originalRun(cmd, args)
	}
	
	// Set up the setup package dependencies
	setup.RunStyledForm = RunStyledForm
	
	// Set up theme-based text functions for setup package
	setup.TitleText = TitleText
	setup.ErrorText = ErrorText
	setup.SuccessText = SuccessText
	setup.InfoTextTitle = InfoTextTitle
	setup.HighlightText = HighlightText
	setup.NormalText = NormalText
	setup.WarningText = WarningText

	// Set up theme-based text functions for presets package
	presets.TitleText = TitleText
	presets.ErrorText = ErrorText
	presets.SuccessText = SuccessText
	presets.InfoTextTitle = InfoTextTitle
	presets.HighlightText = HighlightText
	presets.NormalText = NormalText
	presets.WarningText = WarningText

	// Set up form execution functions for presets package
	presets.RunStyledForm = RunStyledForm
	presets.ShowCancellationMessage = ShowCancellationMessage
	presets.ShowConfirmationCancelledMessage = ShowConfirmationCancelledMessage

	// Set up UI functions for presets package with type adapters
	presets.InputPresetName = inputPresetName
	presets.SelectPlatform = selectPlatform
	presets.SelectParameters = selectParameters
	presets.BuildPresetCommand = func(preset *setup.Preset, config *setup.Config) *presets.CommandInfo {
		// Convert to main package types and call original function
		mainPreset := (*Preset)(preset)
		mainConfig := (*Config)(config)
		cmdInfo := buildPresetCommand(mainPreset, mainConfig)
		// Convert back to presets package type
		return &presets.CommandInfo{
			BaseCommand: cmdInfo.BaseCommand,
			EnvVars:     cmdInfo.EnvVars,
			WorkingDir:  cmdInfo.WorkingDir,
		}
	}
	presets.ExecuteCommand = func(cmdInfo *presets.CommandInfo) {
		// Convert to main package type and call original function
		mainCmdInfo := &CommandInfo{
			BaseCommand: cmdInfo.BaseCommand,
			EnvVars:     cmdInfo.EnvVars,
			WorkingDir:  cmdInfo.WorkingDir,
		}
		executeCommand(mainCmdInfo)
	}
	presets.ShowMoreMenu = func(config *setup.Config) {
		// Convert to main package type and call original function
		mainConfig := (*Config)(config)
		showMoreMenu(mainConfig)
	}

	// Set up theme-based text functions for parameters package
	parameters.TitleText = TitleText
	parameters.ErrorText = ErrorText
	parameters.SuccessText = SuccessText
	parameters.InfoTextTitle = InfoTextTitle
	parameters.HighlightText = HighlightText
	parameters.NormalText = NormalText
	parameters.WarningText = WarningText

	// Set up form execution functions for parameters package
	parameters.RunStyledForm = RunStyledForm
	parameters.ShowCancellationMessage = ShowCancellationMessage
	parameters.ShowConfirmationCancelledMessage = ShowConfirmationCancelledMessage

	// Set up navigation functions for parameters package
	parameters.ShowMoreMenu = func(config *setup.Config) {
		// Convert to main package type and call original function
		mainConfig := (*Config)(config)
		showMoreMenu(mainConfig)
	}
	
	// Add the setup command
	rootCmd.AddCommand(setup.SetupCmd)
	
	// Disable completion
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.CompletionOptions.DisableDescriptions = true
}

func main() {
	// Pass config path to setup package
	if configPath != "" {
		setup.SetConfigPath(configPath)
	}
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}