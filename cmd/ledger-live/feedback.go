package main

import (
	"fmt"
	"ledger-live-starter/cmd/ledger-live/ui"
)

const (
	// GitHub repository URL for feedback
	GitHubURL = "https://github.com/philipptpunkt/ledger-live-starter"
)

// ShowCancellationMessage displays a nicely formatted cancellation message with feedback encouragement
func ShowCancellationMessage() {
	content := fmt.Sprintf(
		"%s\n\n%s\n%s",
		"❌ Selection cancelled",
		"💡 Have feedback or suggestions?",
		"Visit: "+GitHubURL,
	)
	
	fmt.Println("\n" + ui.CreateBoxWrapper(content))
}

// ShowGoodbyeMessage displays a nicely formatted goodbye message 
func ShowGoodbyeMessage() {
	content := fmt.Sprintf(
		"%s\n\n%s\n%s\n\n%s",
		"✓ Goodbye!",
		"Thanks for using Ledger Live Starter!",
		"💡 Feedback? Issues? Ideas?",
		"Visit: "+GitHubURL,
	)
	
	fmt.Println("\n" + ui.CreateBoxWrapper(content))
}

// ShowConfirmationCancelledMessage displays a cancellation message specifically for confirmation dialogs
func ShowConfirmationCancelledMessage() {
	content := fmt.Sprintf(
		"%s\n\n%s\n%s",
		"❌ Confirmation cancelled",
		"💡 Have feedback or suggestions?",
		"Visit: "+GitHubURL,
	)
	
	fmt.Println("\n" + ui.CreateBoxWrapper(content))
}
