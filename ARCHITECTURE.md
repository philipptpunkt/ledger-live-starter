# Project Architecture

This document describes the organizational structure of the ledger-live-starter project.

## Directory Structure

```
ledger-live-starter/
├── cmd/
│   └── ledger-live/           # Main application source code
│       ├── main.go            # Entry point and CLI setup
│       ├── start.go           # Start command and main menu
│       ├── manual.go          # Manual start flow and config loading
│       ├── config.go          # Setup command and configuration
│       ├── presets.go         # Preset management functions
│       ├── parameters.go      # Parameter management functions
│       ├── shared.go          # Shared UI components
│       ├── colors.go          # Color and styling utilities
│       └── version.go         # Version command
├── scripts/                   # Installation scripts
│   ├── install.sh             # Unix/macOS installer
│   └── install.ps1            # Windows PowerShell installer
├── .github/
│   └── workflows/
│       └── release.yml        # GitHub Actions for releases
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
├── Makefile                   # Build and development commands
├── README.md                  # User documentation
├── ARCHITECTURE.md            # This file
└── .gitignore                 # Git ignore patterns
```

## File Descriptions

### Core Application (`cmd/ledger-live/`)

- **`main.go`**: Application entry point, sets up Cobra CLI framework and global flags
- **`start.go`**: Implements the main `start` command, displays interactive menus
- **`manual.go`**: Handles manual start flow, config loading, and path management
- **`config.go`**: Implements setup command and initial configuration wizard
- **`presets.go`**: All preset-related functionality (create, edit, delete, run)
- **`parameters.go`**: Parameter management (add, edit, delete, show)
- **`shared.go`**: Reusable UI components for platform/parameter selection
- **`colors.go`**: ANSI color codes and text styling utilities
- **`version.go`**: Version information command with build details

### Installation Scripts (`scripts/`)

- **`install.sh`**: Bash script for Unix/macOS installation
- **`install.ps1`**: PowerShell script for Windows installation

### Build and CI/CD

- **`Makefile`**: Development and build commands
- **`.github/workflows/release.yml`**: Automated cross-platform builds and releases

## Design Patterns

### Command Structure

Uses the standard Go CLI pattern with Cobra:

- Root command with global flags
- Subcommands for different operations
- Persistent flags for configuration

### Configuration Management

- Default location: `~/.ledger-live/config.json`
- Customizable via `--config` flag or `LEDGER_LIVE_STARTER_CONFIG` environment variable
- Auto-setup on first run if config doesn't exist

### UI Architecture

- `huh` library for interactive forms and menus
- Shared UI components in `shared.go`
- Consistent color scheme via `colors.go`

### Error Handling

- Graceful fallbacks when config is missing
- User-friendly error messages
- Setup mode triggers automatically

## Build System

### Local Development

```bash
make dev           # Run without building
make build         # Build single binary
make build-all     # Cross-platform builds
```

### Release Process

1. Tag release: `git tag v1.0.0 && git push origin v1.0.0`
2. GitHub Actions automatically builds for all platforms
3. Releases published to GitHub with binaries

## Installation Flow

1. User runs install script from GitHub
2. Script detects OS/architecture
3. Downloads appropriate binary from latest GitHub release
4. Installs to `~/.ledger-live/`
5. Adds to PATH
6. On first run, tool triggers setup mode if no config exists

## Configuration Flow

1. Tool checks for config at determined path
2. If missing, runs setup mode automatically
3. Setup creates config directory and prompts for Ledger Live path
4. Default parameters are loaded and user can add custom ones
5. Config saved and tool ready for use
