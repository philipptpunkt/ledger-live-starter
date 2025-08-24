# Ledger Live Starter

🚀 A CLI tool to help users start the ledger-live app with ease.

## Features

- **One-Command Install** - Install directly from GitHub
- **Auto-Setup** - Automatically runs setup on first use
- **Interactive Menu** - Beautiful terminal UI with colors and styling
- **Multiple Platforms** - Support for Mobile and Desktop
- **Custom Presets** - Create and manage your own presets
- **Parameter Management** - Add, edit, and delete custom parameters
- **Global Installation** - Install once, use anywhere
- **Cross-Platform** - Works on macOS, Linux, and Windows
- **Self-Contained** - Everything in `~/.ledger-live/` directory

## Quick Install

### Unix/macOS (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/philipptpunkt/ledger-live-starter/refs/heads/main/scripts/install.sh | bash
```

> **⚠️ Important**: After installation, **restart your terminal** (or run `source ~/.bashrc` / `source ~/.zshrc`) for the PATH changes to take effect.

### Windows (PowerShell as Administrator)

```powershell
iwr -useb https://raw.githubusercontent.com/philipptpunkt/ledger-live-starter/refs/heads/main/scripts/install.ps1 | iex
```

> **⚠️ Important**: After installation, **restart your terminal** (or open a new terminal window) for the PATH changes to take effect. The `ledger-live` command won't be available until you do this.

### Manual Download

Download the latest binary from [Releases](https://github.com/philipptpunkt/ledger-live-starter/releases)

## First Run

```bash
ledger-live start  # Automatically runs setup on first use
```

## Alternative Installation Methods

### Option 1: Install from source (Go required)

```bash
git clone https://github.com/philipptpunkt/ledger-live-starter
cd ledger-live-starter
go install ./cmd/ledger-live
```

### Option 2: Build manually

```bash
# Clone the repository
git clone https://github.com/philipptpunkt/ledger-live-starter
cd ledger-live-starter

# Build the binary
go build -o ledger-live ./cmd/ledger-live

# Move to a directory in your PATH (optional)
sudo mv ledger-live /usr/local/bin/
```

## Usage

### Start Ledger Live

```bash
ledger-live start
```

This opens an interactive menu where you can:

- Select from saved presets
- Start manually with custom parameters
- Access more options (preset/parameter management)

### Run Initial Setup

```bash
ledger-live setup
```

### Use Custom Config File

```bash
ledger-live start --config /path/to/config.json
```

### View Version

```bash
ledger-live version
```

## Configuration

The tool stores configuration in `~/.ledger-live/config.json` by default.

You can customize the config location using:

- `--config` flag: `ledger-live start --config /path/to/config.json`
- Environment variable: `export LEDGER_LIVE_STARTER_CONFIG=/path/to/config.json`

### Sample Configuration

```json
{
  "ledger-live-path": "/Users/username/path/to/ledger-live",
  "parameters": [
    {
      "name": "Skip onboarding",
      "env_var": "SKIP_ONBOARDING=1",
      "description": "Skip the onboarding process"
    },
    {
      "name": "Debug mode",
      "env_var": "DEBUG_MODE=true",
      "description": "Enable debug logging"
    }
  ],
  "presets": [
    {
      "name": "🚀 Mobile Dev",
      "platform": "mobile",
      "parameters": ["Skip onboarding", "Debug mode"]
    },
    {
      "name": "🖥️ Desktop Full",
      "platform": "desktop",
      "parameters": ["Skip onboarding"]
    }
  ]
}
```

## Development

```bash
# Run directly
go run ./cmd/ledger-live

# Run with arguments
go run ./cmd/ledger-live start
go run ./cmd/ledger-live setup

# Build for development (using Makefile)
make build
./ledger-live start

# Or build manually
go build -o ledger-live ./cmd/ledger-live
./ledger-live start
```

## Commands

| Command               | Description                           |
| --------------------- | ------------------------------------- |
| `ledger-live start`   | Interactive menu to start Ledger Live |
| `ledger-live setup`   | Run initial setup or reconfigure      |
| `ledger-live version` | Show version information              |
| `ledger-live --help`  | Show help information                 |

## Directory Structure

After installation, everything is contained in `~/.ledger-live/`:

```
~/.ledger-live/
├── ledger-live          # Binary executable
└── config.json          # Configuration file
```

## Uninstall

To completely remove ledger-live-starter:

```bash
# Remove the installation directory
rm -rf ~/.ledger-live

# Remove from PATH (if you want to remove the PATH entry)
# Edit your shell config file (~/.zshrc, ~/.bashrc, etc.) and remove the line:
# export PATH="$HOME/.ledger-live:$PATH"
```

## Requirements

- Node.js and pnpm (for running Ledger Live)
- Ledger Live repository cloned locally
- Go 1.21+ (only for building from source)

## Updating

To update to the latest version, simply run the install command again:

**Unix/macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/philipptpunkt/ledger-live-starter/refs/heads/main/scripts/install.sh | bash
```

**Windows (PowerShell as Administrator):**

```powershell
iwr -useb https://raw.githubusercontent.com/philipptpunkt/ledger-live-starter/refs/heads/main/scripts/install.ps1 | iex
```

> **💡 Tip**: No need to restart your terminal when updating - the PATH is already configured from the initial installation.

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details on:

- 📝 **Commit Message Convention** - We use Conventional Commits for automatic releases
- 🚀 **Release Process** - Automated with Release Please
- 🔄 **Development Workflow** - From feature to release

## Releases

This project uses [Release Please](https://github.com/googleapis/release-please) for automated releases:

- ✅ **Automatic version bumps** based on commit messages
- ✅ **Generated changelogs** from conventional commits
- ✅ **Cross-platform binaries** built and uploaded automatically
- ✅ **Semantic versioning** following [semver](https://semver.org/)

Latest release: [GitHub Releases](https://github.com/philipptpunkt/ledger-live-starter/releases)
