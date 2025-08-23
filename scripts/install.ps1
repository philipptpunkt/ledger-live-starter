# Ledger Live Starter Installer for Windows
param(
    [string]$InstallDir = "$env:USERPROFILE\.ledger-live"
)

$ErrorActionPreference = "Stop"

# Configuration
$REPO = "philipptpunkt/ledger-live-starter"
$BINARY_NAME = "ledger-live.exe"

Write-Host "Ledger Live Starter Installer" -ForegroundColor Blue -BackgroundColor Black
Write-Host "Installing to: $InstallDir" -ForegroundColor Yellow
Write-Host

# Create installation directory
Write-Host "Creating installation directory..." -ForegroundColor Blue
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

# Get latest release
Write-Host "Fetching latest release..." -ForegroundColor Blue
$releaseUrl = "https://api.github.com/repos/$REPO/releases/latest"
$release = Invoke-RestMethod -Uri $releaseUrl
$downloadUrl = ($release.assets | Where-Object { $_.name -eq "ledger-live-windows-amd64.exe" }).browser_download_url

if (-not $downloadUrl) {
    Write-Host "Error: Could not find Windows binary" -ForegroundColor Red
    Write-Host "Available releases: https://github.com/$REPO/releases" -ForegroundColor Yellow
    exit 1
}

Write-Host "Download URL: $downloadUrl" -ForegroundColor Green

# Download binary
Write-Host "Downloading ledger-live..." -ForegroundColor Blue
$binaryPath = "$InstallDir\$BINARY_NAME"
Invoke-WebRequest -Uri $downloadUrl -OutFile $binaryPath

# Add to PATH
Write-Host "Adding to PATH..." -ForegroundColor Blue
$currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($currentPath -notlike "*$InstallDir*") {
    $newPath = "$InstallDir;$currentPath"
    [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    $env:Path = "$InstallDir;$env:Path"
    Write-Host "Added $InstallDir to PATH" -ForegroundColor Green
} else {
    Write-Host "$InstallDir already in PATH" -ForegroundColor Yellow
}

# Verify installation
Write-Host "Verifying installation..." -ForegroundColor Blue
try {
    & $binaryPath version | Out-Null
    Write-Host "✓ Installation successful!" -ForegroundColor Green
} catch {
    Write-Host "⚠ Installation completed, but verification failed" -ForegroundColor Yellow
    Write-Host "  You may need to restart your terminal" -ForegroundColor Yellow
}

Write-Host
Write-Host "Installation Complete!" -ForegroundColor Green -BackgroundColor Black
Write-Host
Write-Host "Usage:" -ForegroundColor White
Write-Host "  ledger-live start    - Start with interactive menu" -ForegroundColor Green
Write-Host "  ledger-live setup    - Run initial setup" -ForegroundColor Green
Write-Host "  ledger-live --help   - Show help" -ForegroundColor Green
Write-Host
Write-Host "Note: You may need to restart your terminal" -ForegroundColor Yellow
Write-Host "Config location: $InstallDir\config.json" -ForegroundColor Cyan
