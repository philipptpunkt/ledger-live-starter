package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	
	// Disable completion
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.CompletionOptions.DisableDescriptions = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
