package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// CreateBoxWrapper creates a bordered box with gradient border (matching logo colors)
func CreateBoxWrapper(content string) string {
	// Define adaptive text color (same as theme.go)
	textColor := lipgloss.AdaptiveColor{
		Light: "#1a1a1a", // Dark text for light backgrounds
		Dark:  "#ffffff", // White text for dark backgrounds
	}
	
	// Create a basic bordered box with lipgloss for proper alignment
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Foreground(textColor).    // Apply adaptive text color
		Padding(1, 2).
		Width(70).
		Align(lipgloss.Center)
	
	// Render the content with the proper structure
	renderedBox := boxStyle.Render(content)
	
	// Apply brand gradient to border characters only
	return ApplyGradientToBorderOnly(renderedBox, BrandStartColor, BrandEndColor)
}
