# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is `tfren` - a Go CLI tool that organizes Terraform files according to Obay's Terraform Naming Convention (OTN). The tool parses `.tf` files and:
- **Splits** multi-block files into individual files (one block per file)
- **Renames** files to follow the pattern `<block_type>.<provider>.<name>.tf`
- **Validates** file organization against the naming convention

## Architecture

```
tfren/
├── cmd/tfren/main.go          # Entry point
├── internal/
│   ├── cli/                  # Cobra CLI commands
│   │   ├── root.go           # Root command, global flags
│   │   ├── split.go          # Split subcommand
│   │   ├── rename.go         # Rename subcommand
│   │   ├── organize.go       # Organize subcommand (default)
│   │   ├── validate.go       # Validate subcommand
│   │   └── summary.go        # Output helpers
│   ├── config/config.go      # Viper configuration
│   ├── git/status.go         # Git repository checks
│   ├── hcl/                   # HCL parsing
│   │   ├── parser.go         # File parsing
│   │   ├── block.go          # Block structure
│   │   └── comments.go       # Comment handling
│   ├── naming/convention.go  # OTN naming rules
│   ├── operations/           # Core operations
│   │   ├── split.go
│   │   ├── rename.go
│   │   ├── organize.go
│   │   └── validate.go
│   └── output/               # Output formatters
│       ├── console.go        # Colored console output
│       └── json.go           # JSON output for LLMs
└── go.mod
```

## Development Commands

### Building

```bash
go build -o tfren ./cmd/tfren
go build -ldflags="-s -w -X github.com/obay/tfren/internal/cli.Version=<version>" -o tfren ./cmd/tfren
```

### Testing locally

```bash
./tfren --help
./tfren validate                    # Check naming convention
./tfren validate --json             # JSON output for LLMs
./tfren split --dry-run             # Preview split operation
./tfren rename --dry-run            # Preview rename operation
./tfren organize --dry-run          # Preview split+rename
./tfren                             # Run organize (default)
```

### Release Process

- Uses GoReleaser for automated releases to GitHub, Homebrew tap, and Scoop bucket
- GoReleaser config: `.goreleaser.yaml`

### Dependencies

- `github.com/spf13/cobra`: CLI framework
- `github.com/spf13/viper`: Configuration management
- `github.com/hashicorp/hcl/v2`: HCL parsing
- `github.com/charmbracelet/log`: Structured logging
- `github.com/fatih/color`: Terminal colors

## File Naming Convention (OTN)

The tool generates filenames based on Terraform block types:

| Block Type | Pattern | Example |
|------------|---------|---------|
| Resource | `resource.<provider>.<name>.tf` | `resource.aws_instance.web.tf` |
| Data | `data.<provider>.<name>.tf` | `data.aws_ami.ubuntu.tf` |
| Provider | `provider.<name>.tf` | `provider.aws.tf` |
| Provider (alias) | `provider.<name>.<alias>.tf` | `provider.aws.us-east-1.tf` |
| Variable | `variable.<name>.tf` | `variable.environment.tf` |
| Output | `output.<name>.tf` | `output.vpc_id.tf` |
| Module | `module.<name>.tf` | `module.vpc.tf` |
| Locals | `locals.tf` | `locals.tf` |
| Terraform | `terraform.tf` | `terraform.tf` |

## Documentation

- `PRD.md`: Product Requirements Document
- `TDD.md`: Technical Design Document
- `STORIES.md`: Story Map & Epic Breakdown
