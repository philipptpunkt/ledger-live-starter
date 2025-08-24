package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// GitHubRelease represents the response from GitHub's releases API
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
}

// UpdateInfo contains information about available updates
type UpdateInfo struct {
	HasUpdate      bool
	LatestVersion  string
	CurrentVersion string
	UpdateMessage  string
}

// checkForUpdates checks GitHub API for the latest release
func checkForUpdates() *UpdateInfo {
	info := &UpdateInfo{
		CurrentVersion: strings.TrimPrefix(Version, "v"),
		HasUpdate:      false,
	}

	// Skip update check for development versions
	if Version == "dev" || Version == "unknown" {
		return info
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 3 * time.Second, // Quick timeout to avoid blocking startup
	}

	// Make request to GitHub API
	resp, err := client.Get("https://api.github.com/repos/philipptpunkt/ledger-live-starter/releases/latest")
	if err != nil {
		// Silently fail - don't show errors for update checks
		return info
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return info
	}

	// Parse response
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return info
	}

	// Clean version strings for comparison
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := strings.TrimPrefix(Version, "v")

	info.LatestVersion = latestVersion

	// Compare versions
	if isNewerVersion(latestVersion, currentVersion) {
		info.HasUpdate = true
		info.UpdateMessage = createUpdateMessage(latestVersion)
	}

	return info
}

// isNewerVersion compares two semantic version strings
func isNewerVersion(latest, current string) bool {
	// Simple version comparison - split by dots and compare each part
	latestParts := strings.Split(latest, ".")
	currentParts := strings.Split(current, ".")

	// Pad shorter version with zeros
	maxLen := len(latestParts)
	if len(currentParts) > maxLen {
		maxLen = len(currentParts)
	}

	for len(latestParts) < maxLen {
		latestParts = append(latestParts, "0")
	}
	for len(currentParts) < maxLen {
		currentParts = append(currentParts, "0")
	}

	// Compare each part
	for i := 0; i < maxLen; i++ {
		if latestParts[i] > currentParts[i] {
			return true
		} else if latestParts[i] < currentParts[i] {
			return false
		}
	}

	return false
}

// createUpdateMessage creates a styled update notification message
func createUpdateMessage(latestVersion string) string {
	// Use same amber/orange color for both lines
	updateColor := lipgloss.AdaptiveColor{
		Light: "#f59e0b", // Amber for light backgrounds
		Dark:  "#fbbf24", // Lighter amber for dark backgrounds
	}

	// Create styles with same color
	updateStyle := lipgloss.NewStyle().
		Foreground(updateColor).
		Align(lipgloss.Center).
		MarginLeft(4)

	commandStyle := lipgloss.NewStyle().
		Foreground(updateColor). // ← Changed to match updateColor
		Align(lipgloss.Center).
		MarginLeft(4)

	// Determine OS-specific update command
	var updateCommand string
	switch runtime.GOOS {
	case "windows":
		updateCommand = "iwr -useb https://raw.githubusercontent.com/philipptpunkt/ledger-live-starter/refs/heads/main/scripts/install.ps1 | iex"
	default:
		updateCommand = "curl -fsSL https://raw.githubusercontent.com/philipptpunkt/ledger-live-starter/refs/heads/main/scripts/install.sh | bash"
	}

	// Build multi-line message
	var message strings.Builder
	message.WriteString(updateStyle.Render(fmt.Sprintf("🚀 New version v%s available!", latestVersion)))
	message.WriteString("\n")
	message.WriteString(commandStyle.Render("Update with: " + updateCommand))

	return message.String()
}

// getVersionOrUpdateDisplay returns either version info or update notification
func getVersionOrUpdateDisplay() string {
	// Check for updates (with timeout)
	updateInfo := checkForUpdates()

	if updateInfo.HasUpdate {
		return updateInfo.UpdateMessage
	}

	// Fall back to showing current version
	textColor := lipgloss.AdaptiveColor{
		Light: "#1a1a1a", // Dark text for light backgrounds
		Dark:  "#ffffff", // White text for dark backgrounds
	}

	versionStyle := lipgloss.NewStyle().
		Foreground(textColor).
		Align(lipgloss.Center).
		MarginLeft(4)

	return versionStyle.Render(fmt.Sprintf("v%s", Version))
}