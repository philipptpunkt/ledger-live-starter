# Contributing to Ledger Live Starter

Thank you for contributing to ledger-live-starter! This guide will help you understand our development workflow and release process.

## 🚀 Release Process

We use [Release Please](https://github.com/googleapis/release-please) for automated releases and changelog generation. This means:

- **No manual version bumps** - versions are automatically determined from commit messages
- **Automatic changelog generation** - based on your commit messages
- **Automatic releases** - when release PRs are merged

## 📝 Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/) specification. Your commit messages should be structured as follows:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- **`feat:`** - A new feature (triggers **minor** version bump)
- **`fix:`** - A bug fix (triggers **patch** version bump)
- **`docs:`** - Documentation only changes
- **`style:`** - Changes that do not affect the meaning of the code
- **`refactor:`** - A code change that neither fixes a bug nor adds a feature
- **`perf:`** - A code change that improves performance
- **`test:`** - Adding missing tests or correcting existing tests
- **`build:`** - Changes that affect the build system or external dependencies
- **`ci:`** - Changes to our CI configuration files and scripts
- **`chore:`** - Other changes that don't modify src or test files
- **`revert:`** - Reverts a previous commit

### Breaking Changes

To trigger a **major** version bump, include `BREAKING CHANGE:` in the footer or add `!` after the type:

```bash
feat!: redesign CLI interface
# or
feat: new command structure

BREAKING CHANGE: CLI commands have been restructured
```

### Examples

```bash
# Minor version bump (1.0.0 → 1.1.0)
feat: add preset export functionality
feat(presets): add import/export commands

# Patch version bump (1.0.0 → 1.0.1)
fix: resolve environment variable execution issue
fix(commands): properly handle working directory

# No version bump
docs: update installation instructions
style: format code with prettier
chore: update dependencies

# Major version bump (1.0.0 → 2.0.0)
feat!: redesign CLI interface
feat: change command structure

BREAKING CHANGE: All CLI commands have been restructured for better usability
```

## 🔄 Development Workflow

### 1. Development

```bash
# Create feature branch
git checkout -b feature/my-new-feature

# Make changes with conventional commits
git commit -m "feat: add new awesome feature"
git commit -m "fix: handle edge case in feature"
git commit -m "docs: update README with feature info"
```

### 2. Pull Request

- Create PR to `main` branch
- Use conventional commit format in PR title
- All commits should follow conventional format

### 3. Release Process

1. **Merge PR to `main`** → Release Please analyzes commits
2. **Release Please creates release PR** with:
   - Updated CHANGELOG.md
   - Version bump in relevant files
   - Release notes
3. **Review and merge release PR** → Automatic release with binaries

## 🏗️ Local Development

### Setup

```bash
# Clone repository
git clone https://github.com/philipptpunkt/ledger-live-starter
cd ledger-live-starter

# Install dependencies
make setup-dev

# Run in development mode
make dev

# Build locally
make build

# Build for all platforms
make build-all
```

### Testing

```bash
# Run tests
make test

# Test build
make build
./ledger-live version
./ledger-live start
```

## 📋 Release Please Configuration

Our Release Please setup:

- **Release type**: `go`
- **Changelog sections**: Features, Bug Fixes, Performance, Documentation, etc.
- **Version file updates**: `cmd/ledger-live/version.go`
- **Tag format**: `v1.2.3` (includes `v` prefix)

## 🤝 Pull Request Guidelines

1. **Use conventional commit format** in PR title
2. **Keep PRs focused** - one feature/fix per PR
3. **Update documentation** if needed
4. **Test your changes** locally
5. **Follow existing code style**

## 📚 Helpful Resources

- [Conventional Commits Specification](https://www.conventionalcommits.org/)
- [Release Please Documentation](https://github.com/googleapis/release-please)
- [Semantic Versioning](https://semver.org/)

## ❓ Questions?

If you have questions about the contribution process, feel free to:

- Open an issue for discussion
- Ask in your pull request
- Check existing issues and PRs for similar questions
