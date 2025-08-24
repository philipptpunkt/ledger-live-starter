# Project Architecture

This document describes the organizational structure of the ledger-live-starter project.

## Directory Structure

```
ledger-live-starter/
├── cmd/
│   └── ledger-live/                    # Main application source code
│       ├── main.go                     # Entry point and CLI setup
│       ├── start.go                    # Start command and main menu
│       ├── start_manual.go             # Manual start flow and config loading
│       ├── command.go                  # Command building and execution
│       ├── parameters.go               # Parameter management functions
│       ├── shared.go                   # Shared UI components
│       ├── theme.go                    # Theme colors and text styling
│       ├── feedback.go                 # User feedback messages
│       ├── version.go                  # Version command
│       ├── types.go                    # Type aliases for imported packages
│       ├── setup/                      # Configuration and setup package
│       │   ├── setup.go                # Setup command and configuration wizard
│       │   └── config_helpers.go       # Config structures and utilities
│       ├── presets/                    # Preset management package
│       │   ├── shared.go               # Common utilities and dependency injection
│       │   ├── create.go               # Preset creation functionality
│       │   ├── edit.go                 # Preset editing functionality
│       │   ├── delete.go               # Preset deletion functionality
│       │   └── management.go           # Navigation and menu management
│       └── ui/                         # UI components and styling
│           ├── gradient.go             # Gradient color utilities
│           ├── logo.go                 # Application logo rendering
│           └── box_wrapper.go          # Styled box components
├── scripts/                            # Installation scripts
│   ├── install.sh                      # Unix/macOS installer
│   └── install.ps1                     # Windows PowerShell installer
├── .github/
│   └── workflows/
│       └── release.yml                 # GitHub Actions for releases
├── go.mod                              # Go module definition
├── go.sum                              # Go module checksums
├── Makefile                            # Build and development commands
├── README.md                           # User documentation
├── ARCHITECTURE.md                     # This file
└── .gitignore                          # Git ignore patterns
```

## File Descriptions

### Core Application (`cmd/ledger-live/`)

#### Main Files

- **`main.go`**: Application entry point, sets up Cobra CLI framework, dependency injection, and global flags
- **`start.go`**: Implements the main `start` command, displays interactive menus
- **`start_manual.go`**: Handles manual start flow and interactive platform/parameter selection
- **`command.go`**: Command building and execution logic with environment variables
- **`parameters.go`**: Parameter management (add, edit, delete, show)
- **`shared.go`**: Reusable UI components for platform/parameter selection
- **`theme.go`**: Centralized theme system with adaptive colors and text styling functions
- **`feedback.go`**: User feedback messages and notifications
- **`version.go`**: Version information command with build details
- **`types.go`**: Type aliases to avoid redeclaration after package restructuring

#### Setup Package (`setup/`)

- **`setup.go`**: Setup command implementation and configuration wizard
- **`config_helpers.go`**: Configuration structures, file I/O, and utility functions

#### Presets Package (`presets/`)

- **`shared.go`**: Common utilities, config loading, platform conversion, dependency injection
- **`create.go`**: Preset creation logic with shared core functionality (eliminates code duplication)
- **`edit.go`**: Preset editing functionality with form validation
- **`delete.go`**: Preset deletion with confirmation dialogs and bulk operations
- **`management.go`**: Main preset management menu and navigation

#### UI Package (`ui/`)

- **`gradient.go`**: Shared gradient color calculations and text styling
- **`logo.go`**: Application logo rendering with gradients
- **`box_wrapper.go`**: Styled box components with gradient borders

### Installation Scripts (`scripts/`)

- **`install.sh`**: Bash script for Unix/macOS installation
- **`install.ps1`**: PowerShell script for Windows installation

### Build and CI/CD

- **`Makefile`**: Development and build commands
- **`.github/workflows/release.yml`**: Automated cross-platform builds and releases

## Architecture Patterns

### Modular Package Design

The codebase follows a clean modular architecture:

- **`main`**: Entry point with dependency injection
- **`setup`**: Configuration management and initial setup
- **`presets`**: Complete preset lifecycle management
- **`ui`**: Reusable UI components and styling

### Dependency Injection

Functions are injected from `main.go` into subpackages to maintain clean boundaries:

```go
// Theme functions injected into all packages
setup.TitleText = TitleText
presets.ErrorText = ErrorText

// UI functions injected with type adapters
presets.BuildPresetCommand = func(preset *setup.Preset, config *setup.Config) *presets.CommandInfo {
    // Type conversion and delegation
}
```

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
- `lipgloss` for advanced terminal styling and adaptive colors
- Centralized theme system in `theme.go`
- Reusable components in `shared.go` and `ui/` package
- Consistent color scheme with adaptive light/dark mode support

### Error Handling

- Graceful fallbacks when config is missing
- User-friendly error messages with proper color coding
- Setup mode triggers automatically
- Consistent error display across all packages

## Code Quality Features

### Eliminated Code Duplication

The modular refactoring eliminated significant duplication:

- **Preset Creation**: `createPresetCore()` eliminates ~90% duplication between flows
- **Config Loading**: `loadConfigWithDefaults()` centralizes repeated patterns
- **Error Handling**: `saveConfigWithError()` standardizes config save operations
- **Platform Logic**: `convertPlatformToKey()` removes duplicated switch statements
- **UI Components**: Shared gradient logic in `ui/gradient.go`

### Theme System

Centralized theme system provides:

- **Adaptive Colors**: Automatic light/dark mode detection
- **Consistent Styling**: All text uses theme-based functions
- **Semantic Functions**: `TitleText()`, `ErrorText()`, `SuccessText()`, etc.
- **Easy Customization**: Single point of color definition

### Type Safety

- Type adapters handle package boundary conversions
- Proper error propagation between packages
- Clean interfaces with minimal coupling

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

## Performance Benefits

### Code Reduction

- **Before**: 613 lines in single `presets.go` file
- **After**: ~380 lines across 5 focused files in `presets/` package
- **Eliminated**: ~230+ lines of duplicated code (~38% reduction)

### Maintainability Improvements

- Single responsibility principle for each file
- Clear separation of concerns
- Easier testing with smaller, focused functions
- Type-safe package boundaries
- Consistent error handling patterns
