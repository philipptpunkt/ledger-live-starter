package main

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	// Default margin left for consistent spacing across all theme elements
	defaultMarginLeft = 4
)

// GetCustomTheme returns a complete custom theme for huh forms
// Based on the official huh theme structure with our custom colors and styling
func GetCustomTheme() *huh.Theme {
	// Define adaptive colors for consistent theming across light/dark terminals
	primaryColor := lipgloss.AdaptiveColor{
		Light: "#7200c9", // Deep purple for light backgrounds
		Dark:  "#9d4edd", // Lighter purple for dark backgrounds  
	}
	
	titleBlurredColor := lipgloss.AdaptiveColor{
		Light: "#666666", // Dark gray for light backgrounds
		Dark:  "#888888", // Light gray for dark backgrounds
	}
	
	textColor := lipgloss.AdaptiveColor{
		Light: "#1a1a1a", // Dark text for light backgrounds
		Dark:  "#ffffff", // White text for dark backgrounds
	}

	descriptionColor := lipgloss.AdaptiveColor{
		Light: "#666666", // Dark gray for light backgrounds
		Dark:  "#888888", // Light gray for dark backgrounds
	}
	
	highlightColor := lipgloss.AdaptiveColor{
		Light: "#d85500", // Darker orange for light backgrounds (better contrast)
		Dark:  "#f2830c", // Bright orange for dark backgrounds (matches gradient end)
	}
	
	errorColor := lipgloss.AdaptiveColor{
		Light: "#dc2626", // Red for light backgrounds
		Dark:  "#ef4444", // Lighter red for dark backgrounds
	}
	
	successColor := lipgloss.AdaptiveColor{
		Light: "#16a34a", // Green for light backgrounds
		Dark:  "#22c55e", // Lighter green for dark backgrounds
	}
	
	// Create complete custom theme with all possible huh elements
	return &huh.Theme{
		Focused: huh.FieldStyles{
			// Base styles - applied to all form elements
			Base:           lipgloss.NewStyle().BorderForeground(textColor).Border(lipgloss.ThickBorder(), false, false, false, true).BorderForeground(primaryColor),
			Title:          lipgloss.NewStyle().Foreground(primaryColor).Bold(true).MarginLeft(defaultMarginLeft),
			Description:    lipgloss.NewStyle().Foreground(descriptionColor).MarginLeft(defaultMarginLeft),
			ErrorIndicator: lipgloss.NewStyle().Foreground(errorColor).MarginLeft(defaultMarginLeft),
			ErrorMessage:   lipgloss.NewStyle().Foreground(errorColor).MarginLeft(defaultMarginLeft),
			
			// Select field styles (single select dropdowns)
			SelectSelector:   lipgloss.NewStyle().Foreground(primaryColor).MarginLeft(defaultMarginLeft).Bold(true).SetString("▶ "), // Selection indicator (cursor)
			Option:           lipgloss.NewStyle().Foreground(textColor).MarginLeft(defaultMarginLeft).Bold(true),     // All select options
			NextIndicator:    lipgloss.NewStyle().Foreground(highlightColor), // Navigation arrows
			PrevIndicator:    lipgloss.NewStyle().Foreground(highlightColor), // Navigation arrows
			
			// FilePicker styles (file/directory selection)
			Directory: lipgloss.NewStyle().Foreground(textColor),
			File:      lipgloss.NewStyle().Foreground(textColor),
			
			// Multi-select styles (checkboxes, multiple selection)
			MultiSelectSelector: lipgloss.NewStyle().Foreground(primaryColor).MarginLeft(defaultMarginLeft).SetString("▶ "), // Multi-select cursor
			SelectedOption:      lipgloss.NewStyle().Foreground(highlightColor), // Selected items
			SelectedPrefix:      lipgloss.NewStyle().Foreground(successColor).SetString("✓ "), // Checkmark for selected
			UnselectedOption:    lipgloss.NewStyle().Foreground(textColor),     // Unselected items
			UnselectedPrefix:    lipgloss.NewStyle().Foreground(descriptionColor).SetString("✗ "),     // X mark for unselected
			
			// Text input styles (text fields, text areas)
			TextInput: huh.TextInputStyles{
				Cursor:      lipgloss.NewStyle().Foreground(textColor).MarginLeft(defaultMarginLeft),              	  // Blinking cursor in text fields
				Placeholder: lipgloss.NewStyle().Foreground(descriptionColor).MarginLeft(defaultMarginLeft), // Placeholder text
				Prompt:      lipgloss.NewStyle().Foreground(primaryColor).MarginLeft(defaultMarginLeft),       // Input prompt/label
				Text:        lipgloss.NewStyle().Foreground(textColor).MarginLeft(defaultMarginLeft),        // Entered text
			},
			
			// Confirm/Button styles (Yes/No prompts, buttons)
			FocusedButton: lipgloss.NewStyle().
				Foreground(highlightColor).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(highlightColor).
				Padding(0, 2).
				Bold(true),
			BlurredButton: lipgloss.NewStyle().
				Foreground(descriptionColor).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(descriptionColor).
				Padding(0, 2),
			
			// Card/Note styles (info displays, notes)
			Card:      lipgloss.NewStyle().MarginLeft(defaultMarginLeft).BorderForeground(textColor).Padding(1),
			NoteTitle: lipgloss.NewStyle().Foreground(primaryColor).MarginLeft(defaultMarginLeft).Bold(true),
			Next:      lipgloss.NewStyle().Foreground(highlightColor).MarginLeft(defaultMarginLeft).Bold(true),
		},
		
		// Blurred state (when form loses focus)
		Blurred: huh.FieldStyles{
			// Base styles - muted when not focused
			Base:           lipgloss.NewStyle().BorderForeground(textColor).Border(lipgloss.ThickBorder(), false, false, false, true).BorderForeground(titleBlurredColor),
			Title:          lipgloss.NewStyle().Foreground(titleBlurredColor).Bold(true).MarginLeft(defaultMarginLeft),
			Description:    lipgloss.NewStyle().Foreground(textColor).MarginLeft(defaultMarginLeft).MarginBottom(1),
			ErrorIndicator: lipgloss.NewStyle().Foreground(errorColor).MarginLeft(defaultMarginLeft),
			ErrorMessage:   lipgloss.NewStyle().Foreground(errorColor).MarginLeft(defaultMarginLeft),
			
			// Select field styles - muted when blurred
			SelectSelector:   lipgloss.NewStyle().Foreground(titleBlurredColor).MarginLeft(defaultMarginLeft).Bold(true).SetString("▶ "), // Muted selector
			Option:           lipgloss.NewStyle().Foreground(textColor).MarginLeft(defaultMarginLeft).Bold(true),       // Options stay readable
			NextIndicator:    lipgloss.NewStyle().Foreground(titleBlurredColor), // Muted navigation
			PrevIndicator:    lipgloss.NewStyle().Foreground(titleBlurredColor), // Muted navigation
			
			// FilePicker styles - muted when blurred
			Directory: lipgloss.NewStyle().Foreground(textColor),
			File:      lipgloss.NewStyle().Foreground(textColor),
			
			// Multi-select styles - keep selected items visible even when blurred
			MultiSelectSelector: lipgloss.NewStyle().Foreground(titleBlurredColor).MarginLeft(defaultMarginLeft).SetString("▶ "), // Muted cursor
			SelectedOption:      lipgloss.NewStyle().Foreground(highlightColor),     // Keep selected visible
			SelectedPrefix:      lipgloss.NewStyle().Foreground(successColor).SetString("✓ "),     // Keep selected markers visible
			UnselectedOption:    lipgloss.NewStyle().Foreground(textColor),         // Regular text
			UnselectedPrefix:    lipgloss.NewStyle().Foreground(descriptionColor).SetString("✗ "),         // Regular markers
			
			// Text input styles - muted when blurred
			TextInput: huh.TextInputStyles{
				Cursor:      lipgloss.NewStyle().Foreground(titleBlurredColor),               // Muted cursor
				Placeholder: lipgloss.NewStyle().Foreground(titleBlurredColor).MarginLeft(defaultMarginLeft), // Muted placeholder
				Prompt:      lipgloss.NewStyle().Foreground(titleBlurredColor).MarginLeft(defaultMarginLeft), // Muted prompt
				Text:        lipgloss.NewStyle().Foreground(textColor).MarginLeft(defaultMarginLeft),         // Text stays readable
			},
			
			// Confirm/Button styles - muted when blurred
			FocusedButton: lipgloss.NewStyle().
				Foreground(titleBlurredColor).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(titleBlurredColor).
				Padding(0, 2).
				MarginLeft(defaultMarginLeft),
			BlurredButton: lipgloss.NewStyle().
				Foreground(textColor).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(textColor).
				Padding(0, 2).
				MarginLeft(defaultMarginLeft),
			
			// Card/Note styles - muted when blurred
			Card:      lipgloss.NewStyle().MarginLeft(defaultMarginLeft).BorderForeground(textColor).Padding(1),
			NoteTitle: lipgloss.NewStyle().Foreground(titleBlurredColor).MarginLeft(defaultMarginLeft),
			Next:      lipgloss.NewStyle().Foreground(titleBlurredColor).MarginLeft(defaultMarginLeft),
		},
	}
}

// RunStyledForm applies our custom theme to a huh form and runs it
func RunStyledForm(form *huh.Form) error {
	return form.WithTheme(GetCustomTheme()).Run()
}

// Theme-based text styling functions using consistent colors from theme
func NormalText(text string) string {
	textColor := lipgloss.AdaptiveColor{
		Light: "#1a1a1a", // Dark text for light backgrounds
		Dark:  "#ffffff", // White text for dark backgrounds
	}
	return lipgloss.NewStyle().Foreground(textColor).Render(text)
}

func TitleText(text string) string {
	primaryColor := lipgloss.AdaptiveColor{
		Light: "#7200c9", // Deep purple for light backgrounds
		Dark:  "#9d4edd", // Lighter purple for dark backgrounds  
	}
	return lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(text)
}

func HighlightText(text string) string {
	highlightColor := lipgloss.AdaptiveColor{
		Light: "#d85500", // Darker orange for light backgrounds (better contrast)
		Dark:  "#f2830c", // Bright orange for dark backgrounds (matches gradient end)
	}
	return lipgloss.NewStyle().Foreground(highlightColor).Render(text)
}

func SuccessText(text string) string {
	successColor := lipgloss.AdaptiveColor{
		Light: "#16a34a", // Green for light backgrounds
		Dark:  "#22c55e", // Lighter green for dark backgrounds
	}
	return lipgloss.NewStyle().Foreground(successColor).Bold(true).Render(text)
}

func ErrorText(text string) string {
	errorColor := lipgloss.AdaptiveColor{
		Light: "#dc2626", // Red for light backgrounds
		Dark:  "#ef4444", // Lighter red for dark backgrounds
	}
	return lipgloss.NewStyle().Foreground(errorColor).Bold(true).Render(text)
}

// Info text titles use a mid blue color
func InfoTextTitle(text string) string {
	infoTitleColor := lipgloss.AdaptiveColor{
		Light: "#0369a1", // Mid blue for light backgrounds
		Dark:  "#38bdf8", // Lighter blue for dark backgrounds
	}
	return lipgloss.NewStyle().Foreground(infoTitleColor).Render(text)
}

func WarningText(text string) string {
	warningColor := lipgloss.AdaptiveColor{
		Light: "#d97706", // Orange/amber for light backgrounds
		Dark:  "#fbbf24", // Yellow for dark backgrounds
	}
	return lipgloss.NewStyle().Foreground(warningColor).Bold(true).Render(text)
}
