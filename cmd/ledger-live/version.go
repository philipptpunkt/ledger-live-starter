package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"ledger-live-starter/cmd/ledger-live/setup"
)

// Version information, set at build time via ldflags
var (
	Version   = "dev"  // Set via: -ldflags "-X main.Version=v1.0.0"
	BuildTime = "unknown"  // Set via: -ldflags "-X main.BuildTime=2024-01-01T00:00:00Z"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of ledger-live-starter`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s %s\n", TitleText("Ledger Live Starter"), HighlightText(Version))
		fmt.Printf("  %s %s\n", InfoTextTitle("Built:"), NormalText(BuildTime))
		fmt.Printf("  %s %s\n", InfoTextTitle("Go:"), NormalText(runtime.Version()))
		fmt.Printf("  %s %s/%s\n", InfoTextTitle("Platform:"), NormalText(runtime.GOOS), NormalText(runtime.GOARCH))
		fmt.Printf("  %s %s\n", InfoTextTitle("Config:"), NormalText(setup.GetConfigPath()))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
