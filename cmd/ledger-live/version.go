package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
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
		fmt.Printf("%s %s\n", BoldColorText("Ledger Live Starter", Cyan), ColorText("v"+Version, Green))
		fmt.Printf("  %s %s\n", ColorText("Built:", Blue), BuildTime)
		fmt.Printf("  %s %s\n", ColorText("Go:", Blue), runtime.Version())
		fmt.Printf("  %s %s/%s\n", ColorText("Platform:", Blue), runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  %s %s\n", ColorText("Config:", Blue), getConfigPath())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
