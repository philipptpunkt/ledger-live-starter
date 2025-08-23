package main

// ANSI color codes and styling
const (
	// Colors
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	White  = "\033[37m"
	
	// Styles
	Bold      = "\033[1m"
	Underline = "\033[4m"
	Reset     = "\033[0m"
	
	// Combined styles for common use cases
	BoldRed    = Bold + Red
	BoldGreen  = Bold + Green
	BoldBlue   = Bold + Blue
	BoldCyan   = Bold + Cyan
	BoldYellow = Bold + Yellow
	BoldPurple = Bold + Purple
)

// Helper functions for styling text
func StyleText(text, style string) string {
	return style + text + Reset
}

func BoldText(text string) string {
	return StyleText(text, Bold)
}

func ColorText(text, color string) string {
	return StyleText(text, color)
}

func BoldColorText(text, color string) string {
	return StyleText(text, Bold+color)
}
