package ui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

// Brand gradient colors - from deep purple to orange
var (
	BrandStartColor = [3]int{114, 0, 201}  // Deep Purple #7200c9
	BrandEndColor   = [3]int{242, 131, 12} // Orange #f2830c
)

// GetGradientColor calculates a specific color from the gradient at a given progress (0.0 to 1.0)
func GetGradientColor(progress float64, startColor [3]int, endColor [3]int) string {
	r := int(float64(startColor[0]) + progress*float64(endColor[0]-startColor[0]))
	g := int(float64(startColor[1]) + progress*float64(endColor[1]-startColor[1]))
	b := int(float64(startColor[2]) + progress*float64(endColor[2]-startColor[2]))
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

// ApplyGradientToText applies a left-to-right gradient to all characters in text
func ApplyGradientToText(text string, startColor [3]int, endColor [3]int) string {
	lines := []string{}
	for _, line := range splitLines(text) {
		if len(line) == 0 {
			lines = append(lines, "")
			continue
		}
		
		gradientLine := ""
		for i, char := range line {
			// Calculate gradient position (0.0 to 1.0)
			progress := float64(i) / float64(len(line)-1)
			if len(line) == 1 {
				progress = 0.5 // Single character gets middle color
			}
			
			// Apply gradient color to character
			gradientColor := GetGradientColor(progress, startColor, endColor)
			colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(gradientColor))
			gradientLine += colorStyle.Render(string(char))
		}
		lines = append(lines, gradientLine)
	}
	
	return joinLines(lines)
}

// ApplyGradientToBorderOnly applies gradient only to border characters, leaving content unchanged
func ApplyGradientToBorderOnly(boxContent string, startColor [3]int, endColor [3]int) string {
	lines := splitLines(boxContent)
	if len(lines) == 0 {
		return boxContent
	}
	
	resultLines := []string{}
	for _, line := range lines {
		if len(line) == 0 {
			resultLines = append(resultLines, "")
			continue
		}
		
		gradientLine := ""
		for j, char := range line {
			charStr := string(char)
			
			// Check if this character is a border character
			isBorderChar := isBorderCharacter(charStr)
			
			if isBorderChar {
				// Calculate gradient position for this character
				progress := float64(j) / float64(len(line)-1)
				if len(line) == 1 {
					progress = 0.5
				}
				
				// Apply gradient color
				gradientColor := GetGradientColor(progress, startColor, endColor)
				colorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(gradientColor))
				gradientLine += colorStyle.Render(charStr)
			} else {
				// Keep content characters unchanged
				gradientLine += charStr
			}
		}
		resultLines = append(resultLines, gradientLine)
	}
	
	return joinLines(resultLines)
}

// Helper function to check if a character is a border character
func isBorderCharacter(char string) bool {
	return char == "╭" || char == "╮" || char == "╯" || char == "╰" || // rounded corners
		char == "─" || char == "│" || // horizontal and vertical lines  
		char == "┌" || char == "┐" || char == "└" || char == "┘" // square corners
}

// Helper function to split text into lines
func splitLines(text string) []string {
	if text == "" {
		return []string{}
	}
	
	lines := []string{}
	current := ""
	for _, char := range text {
		if char == '\n' {
			lines = append(lines, current)
			current = ""
		} else {
			current += string(char)
		}
	}
	lines = append(lines, current)
	return lines
}

// Helper function to join lines back into text
func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}
