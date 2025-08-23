# Setting Up Your First Release

This document explains how to set up Release Please for the first time.

## 🚀 Initial Setup Steps

### 1. Push to Main Branch

First, push all your current changes to the `main` branch:

```bash
# Add all files
git add .

# Use conventional commit format
git commit -m "feat: initial release of ledger-live-starter

- Interactive CLI with preset management
- Cross-platform support (macOS, Linux, Windows)
- Auto-setup mode for first-time users
- Parameter management system
- Colored terminal output
- GitHub-based installation system"

# Push to main
git push origin main
```

### 2. First Release Please Run

After pushing to `main`, Release Please will:

1. **Analyze your commits** since the beginning of the repository
2. **Create a release PR** with:
   - Version 1.0.0 (since this is the first release)
   - Generated CHANGELOG.md with all features
   - Updated version in `cmd/ledger-live/version.go`

### 3. Review and Merge Release PR

1. **Check the generated release PR** (should appear automatically)
2. **Review the changelog** - make sure it looks good
3. **Merge the release PR**
4. **Watch the magic happen**:
   - Tag `v1.0.0` gets created
   - Cross-platform binaries are built automatically
   - GitHub release is created with binaries attached
   - Install scripts can now download the release

## 🔮 Future Releases

After the first release, the workflow is:

```bash
# Make changes with conventional commits
git commit -m "feat: add new awesome feature"
git commit -m "fix: resolve some bug"

# Push to main
git push origin main

# Release Please automatically creates release PR
# Review and merge → automatic release!
```

## 🏷️ Version Bumping

Release Please determines version bumps from commit messages:

- `feat:` → **Minor** (1.0.0 → 1.1.0)
- `fix:` → **Patch** (1.0.0 → 1.0.1)
- `feat!:` or `BREAKING CHANGE:` → **Major** (1.0.0 → 2.0.0)

## 📋 Troubleshooting

### No Release PR Created?

- Check that commits use conventional format (`feat:`, `fix:`, etc.)
- Ensure you're pushing to `main` branch
- Look at GitHub Actions logs for any errors

### Want to Skip a Release?

- Use commit types that don't trigger releases: `docs:`, `chore:`, `style:`, etc.
- Or add `[skip ci]` to commit message

### Manual Release?

If needed, you can manually create a release:

```bash
# Create and push a tag
git tag v1.0.0
git push origin v1.0.0

# GitHub Actions will build and release automatically
```

## 🎉 You're Ready!

Once you push your first conventional commit to `main`, Release Please will take care of the rest!

Delete this file after your first successful release.
