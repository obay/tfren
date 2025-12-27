# tfren

[![Go Version](https://img.shields.io/github/go-mod/go-version/obay/tfren)](https://go.dev/)
[![Release](https://img.shields.io/github/v/release/obay/tfren)](https://github.com/obay/tfren/releases/latest)
[![License](https://img.shields.io/github/license/obay/tfren)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/obay/tfren)](https://goreportcard.com/report/github.com/obay/tfren)

**Terraform file organizer** — Automatically split and rename your Terraform files to follow a consistent naming convention.

## What it does

tfren organizes messy Terraform configurations into clean, single-block files following **Obay's Terraform Naming Convention (OTN)**:

```
Before:                          After:
main.tf (100+ lines)      →      resource.aws_instance.web.tf
                                 resource.aws_security_group.web.tf
                                 variable.instance_type.tf
                                 variable.environment.tf
                                 output.instance_ip.tf
                                 provider.aws.tf
                                 terraform.tf
```

## Features

- **Split** multi-block files into individual files (one block per file)
- **Rename** files to follow the naming pattern `<block_type>.<provider>.<name>.tf`
- **Validate** existing files against the naming convention
- **Recursive** processing of subdirectories
- **Dry-run** mode to preview changes before applying
- **Git-aware** — warns about uncommitted changes and prompts before modifying uncommitted `.tf` files

## Installation

### Homebrew (macOS & Linux)

```bash
brew install obay/tap/tfren
```

### Scoop (Windows)

```powershell
scoop bucket add obay https://github.com/obay/scoop-bucket.git
scoop install tfren
```

### Go Install

```bash
go install github.com/obay/tfren/cmd/tfren@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/obay/tfren/releases/latest).

## Usage

### Organize (default)

Split multi-block files and rename all files to follow the naming convention:

```bash
tfren                    # Organize current directory
tfren -r                 # Organize recursively
tfren -n                 # Dry-run (preview changes)
```

### Split Only

Split multi-block files without renaming:

```bash
tfren split
tfren split -r           # Recursive
```

### Rename Only

Rename files without splitting:

```bash
tfren rename
tfren rename -r          # Recursive
```

### Validate

Check if files follow the naming convention:

```bash
tfren validate
tfren validate --json    # JSON output
```

## Naming Convention

| Block Type | Pattern | Example |
|------------|---------|---------|
| Resource | `resource.<provider>.<name>.tf` | `resource.aws_instance.web.tf` |
| Data | `data.<provider>.<name>.tf` | `data.aws_ami.ubuntu.tf` |
| Variable | `variable.<name>.tf` | `variable.environment.tf` |
| Output | `output.<name>.tf` | `output.vpc_id.tf` |
| Module | `module.<name>.tf` | `module.vpc.tf` |
| Provider | `provider.<name>.tf` | `provider.aws.tf` |
| Provider (alias) | `provider.<name>.<alias>.tf` | `provider.aws.us-east-1.tf` |
| Locals | `locals.tf` | `locals.tf` |
| Terraform | `terraform.tf` | `terraform.tf` |

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--directory` | `-d` | Target directory (default: current) |
| `--recursive` | `-r` | Process subdirectories |
| `--dry-run` | `-n` | Preview changes without applying |
| `--force` | `-f` | Skip confirmation prompts |
| `--no-git-check` | | Skip git status check |
| `--quiet` | `-q` | Suppress non-error output |
| `--json` | | Output in JSON format |
| `--verbose` | `-v` | Verbose output |

## License

MIT License - see [LICENSE](LICENSE) for details.
