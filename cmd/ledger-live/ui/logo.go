package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// GetLogo returns the ASCII logo with purple-to-orange gradient
func GetLogo() string {
	logo := `
╔═════════════════════════════╗
║                             ║
║    Ledger Live Starter      ║
║                             ║
╚═════════════════════════════╝`
	
	// Apply brand gradient to the entire logo
	gradientLogo := ApplyGradientToText(logo, BrandStartColor, BrandEndColor)
	
	// Add left margin for better positioning
	logoWithMargin := lipgloss.NewStyle().
		MarginLeft(4).
		Render(gradientLogo)
	
	return logoWithMargin
}
