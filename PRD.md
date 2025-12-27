# Product Requirements Document (PRD)

## tfrn - Terraform File Organizer

**Version:** 1.0
**Date:** December 2024
**Author:** Obay
**Status:** Draft

---

## 1. Executive Summary

### 1.1 Overview

`tfrn` (Terraform Rename) is a CLI tool for organizing Terraform configuration files according to Obay's Terraform Naming Convention (OTN). It automates file splitting, renaming, and validation to maintain a clean, navigable Terraform codebase.

### 1.2 Problem Statement

Terraform projects often suffer from poor file organization:
- Multiple resources defined in a single `main.tf` file
- Inconsistent file naming that doesn't reflect contents
- Difficulty navigating large Terraform codebases
- Manual effort required to reorganize existing projects

### 1.3 Solution

`tfrn` provides a comprehensive tool that:
1. **Splits** multi-block Terraform files into individual files (one block per file)
2. **Renames** files according to a standardized naming convention
3. **Validates** file organization against the naming convention
4. Provides **machine-readable output** for integration with LLMs and automation tools

---

## 2. Goals and Objectives

### 2.1 Primary Goals

| Goal | Description | Success Metric |
|------|-------------|----------------|
| Completeness | Single tool for all Terraform file organization needs | Split, rename, validate in one binary |
| Usability | Intuitive CLI with subcommands | New users productive in < 5 minutes |
| LLM Integration | Machine-readable output for AI assistants | JSON output mode, clear exit codes |

### 2.2 Non-Goals

- Modifying Terraform resource configurations (only file organization)
- Validating Terraform syntax beyond HCL parsing
- Managing Terraform state or remote backends
- Formatting Terraform code (use `terraform fmt` for that)

---

## 3. Naming Convention (OTN)

The tool enforces **Obay's Terraform Naming Convention (OTN)** from [obay.cloud/post/terraform-naming-convention](https://www.obay.cloud/post/terraform-naming-convention/).

### 3.1 Background

This naming convention is not new—it applies a well-established software engineering practice to Terraform. In languages like C++ and C#, the "one class per file" convention has been standard practice for decades. Files are named after the class they contain, and folder structures mirror namespaces. This approach has proven benefits:

- **Navigation**: Developers can locate code instantly by name
- **Maintainability**: Changes are isolated to specific, predictable files
- **Code Review**: File diffs clearly show what component changed
- **Tooling**: IDEs and search tools work more effectively

OTN brings these same benefits to Terraform by treating each HCL block (resource, data source, variable, etc.) as analogous to a class—one block per file, named to reflect its contents.

### 3.2 Core Principle

> "One block per file, named according to the block type and identifier."

### 3.3 File Naming Patterns

| Block Type | Pattern | Example |
|------------|---------|---------|
| Resource | `resource.<provider_resource>.<name>.tf` | `resource.aws_instance.web-server.tf` |
| Data Source | `data.<provider_resource>.<name>.tf` | `data.aws_ami.ubuntu.tf` |
| Provider | `provider.<name>.tf` | `provider.aws.tf` |
| Provider (alias) | `provider.<name>.<alias>.tf` | `provider.aws.us-east-1.tf` |
| Variable | `variable.<name>.tf` | `variable.environment.tf` |
| Output | `output.<name>.tf` | `output.public_ip.tf` |
| Module | `module.<name>.tf` | `module.vpc.tf` |
| Locals | `locals.tf` | `locals.tf` |
| Terraform | `terraform.tf` | `terraform.tf` |

---

## 4. Functional Requirements

### 4.1 Command Structure

```
tfrn [command] [flags]

Commands:
  split       Split multi-block files into individual files
  rename      Rename files according to naming convention
  organize    Split AND rename in one operation (default)
  validate    Check files against naming convention (no changes)
  help        Help about any command

Global Flags:
  -d, --directory string   Target directory (default: current directory)
  -r, --recursive          Process subdirectories recursively
  -n, --dry-run            Show what would be done without making changes
  -v, --verbose            Verbose output
  -q, --quiet              Suppress non-error output
      --json               Output in JSON format (for LLM/automation)
      --version            Print version information
  -h, --help               Help for tfrn
```

### 4.2 Safety Checks

Before modifying any files, `tfrn` checks whether the target directory is under Git version control:

- **Not a Git repository**: Display a warning that changes cannot be easily undone and recommend initializing a Git repository
- **Git repository with uncommitted changes**: Display a warning that there are uncommitted changes and recommend committing or stashing before proceeding
- **Clean Git repository**: Proceed without warning

These warnings are informational only and do not block the operation. Users can suppress them with `--quiet`.

### 4.3 Commands

#### 4.3.1 `tfrn split`

Split files containing multiple blocks into separate files.

```bash
tfrn split [flags]

Flags:
  --keep-original   Keep the original file after splitting
  --backup          Create .bak backup before splitting
```

**Behavior:**
1. Parse all `.tf` files in target directory
2. Identify files with multiple HCL blocks
3. Extract each block into a separate file
4. Name new files according to OTN convention
5. Remove original file (unless `--keep-original`)

**Comment Handling:**
Comments immediately preceding a block are associated with that block and included in the split file. File-level headers (comments at the very top before any blocks) and "floating" comments (comments not directly attached to a block) are discarded. Inline comments within blocks are always preserved.

#### 4.3.2 `tfrn rename`

Rename single-block files to match naming convention.

```bash
tfrn rename [flags]

Flags:
  --backup   Create .bak backup before renaming
```

**Behavior:**
1. Parse all `.tf` files in target directory
2. Skip files with multiple blocks (warn user)
3. Generate correct filename based on block content
4. Rename file if name doesn't match convention

#### 4.3.3 `tfrn organize`

Combines split and rename operations (default command).

```bash
tfrn organize [flags]
# or simply:
tfrn [flags]
```

**Behavior:**
1. First, split all multi-block files
2. Then, rename all files to match convention
3. Report summary of changes

#### 4.3.4 `tfrn validate`

Check files against naming convention without making changes.

```bash
tfrn validate [flags]

Flags:
  --strict   Exit with error code if any violations found
```

**Behavior:**
1. Parse all `.tf` files
2. Check each file against naming convention
3. Report violations
4. Exit code: 0 = compliant, 1 = violations found (with `--strict`)

---

## 5. Technical Requirements

### 5.1 Technology Stack

| Component | Technology | Rationale |
|-----------|------------|-----------|
| Language | Go 1.25+ | Performance, single binary distribution |
| CLI Framework | [Cobra](https://github.com/spf13/cobra) | Industry standard, subcommand support |
| Configuration | [Viper](https://github.com/spf13/viper) | Config files, env vars, flag binding |
| HCL Parsing | [hashicorp/hcl/v2](https://github.com/hashicorp/hcl) | Official Terraform HCL parser |
| Logging | [charmbracelet/log](https://github.com/charmbracelet/log) | Beautiful, structured logging |
| Colors | [fatih/color](https://github.com/fatih/color) | Cross-platform terminal colors |

### 5.2 Architecture

```
tfrn/
├── cmd/
│   ├── root.go           # Root command, global flags
│   ├── split.go          # Split subcommand
│   ├── rename.go         # Rename subcommand
│   ├── organize.go       # Organize subcommand (default)
│   └── validate.go       # Validate subcommand
├── internal/
│   ├── parser/
│   │   └── hcl.go        # HCL parsing utilities
│   ├── naming/
│   │   └── convention.go # OTN naming logic
│   ├── operations/
│   │   ├── split.go      # File splitting logic
│   │   ├── rename.go     # File renaming logic
│   │   └── validate.go   # Validation logic
│   └── output/
│       ├── console.go    # Human-readable output
│       └── json.go       # JSON output for LLMs
├── main.go               # Entry point
├── version.go            # Version information
└── go.mod
```

### 5.3 Configuration File Support

`tfrn` supports configuration via `.tfrn.json`:

```json
{
  "recursive": false,
  "dry_run": false,
  "verbose": false,
  "output_format": "console",
  "patterns": {
    "resource": "resource.{provider}.{name}.tf",
    "data": "data.{provider}.{name}.tf"
  }
}
```

Configuration priority (highest to lowest):
1. Command-line flags
2. Environment variables (`TFRN_*`)
3. Configuration file in current directory
4. Configuration file in home directory
5. Default values

---

## 6. LLM Integration Requirements

### 6.1 Machine-Readable Output

When `--json` flag is used, output structured JSON for LLM consumption:

```json
{
  "version": "1.0.0",
  "command": "organize",
  "directory": "/path/to/terraform",
  "dry_run": false,
  "results": {
    "files_processed": 15,
    "files_split": 3,
    "files_renamed": 8,
    "files_skipped": 2,
    "files_already_compliant": 2,
    "errors": []
  },
  "changes": [
    {
      "action": "split",
      "source": "main.tf",
      "targets": [
        "resource.aws_instance.web.tf",
        "resource.aws_security_group.web.tf",
        "data.aws_ami.ubuntu.tf"
      ]
    },
    {
      "action": "rename",
      "source": "networking.tf",
      "target": "resource.aws_vpc.main.tf"
    }
  ],
  "violations": [],
  "exit_code": 0
}
```

### 6.2 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success, all operations completed |
| 1 | Partial success, some files skipped |
| 2 | Validation failed (with `--strict`) |
| 3 | Fatal error (parsing, I/O, etc.) |

### 6.3 LLM Usage Examples

**For Warp Terminal AI / Claude Code:**

```bash
# Show what would be organized
tfrn --dry-run --json

# Organize and report results
tfrn organize --json

# Validate before PR
tfrn validate --strict --json
```

**Prompt example for LLMs:**

> "Run `tfrn validate --json` to check if Terraform files follow the naming convention. Parse the JSON output to identify files that need attention."

---

## 7. User Stories

### 7.1 Developer Stories

| ID | As a... | I want to... | So that... |
|----|---------|--------------|------------|
| US-1 | Developer | Split my main.tf into individual files | I can find resources easily |
| US-2 | Developer | Rename files to match convention | My project is consistently organized |
| US-3 | Developer | Preview changes before applying | I don't accidentally break anything |
| US-4 | Developer | Validate files in CI/CD | PRs maintain naming standards |

### 7.2 LLM/Automation Stories

| ID | As a... | I want to... | So that... |
|----|---------|--------------|------------|
| US-5 | LLM (Claude Code) | Get JSON output from tfrn | I can parse and act on results |
| US-6 | CI Pipeline | Run validation with strict mode | Builds fail on naming violations |
| US-7 | Automation script | Process results programmatically | I can integrate with other tools |

---

## 8. Acceptance Criteria

### 8.1 Core Functionality

- [ ] `tfrn split` correctly extracts all block types
- [ ] `tfrn rename` handles all OTN naming patterns
- [ ] `tfrn organize` combines both operations seamlessly
- [ ] `tfrn validate` detects all naming violations
- [ ] `--dry-run` makes no file system changes
- [ ] `--json` produces valid, parseable JSON
- [ ] `--recursive` processes subdirectories correctly

### 8.2 Edge Cases

- [ ] Handles files with syntax errors gracefully
- [ ] Detects and reports file name conflicts
- [ ] Preserves file permissions after rename
- [ ] Works with symlinks appropriately
- [ ] Handles empty directories
- [ ] Handles read-only files with clear error

### 8.3 LLM Integration

- [ ] JSON output is consistent and documented
- [ ] Exit codes are reliable and documented
- [ ] Error messages are actionable
- [ ] Dry-run mode provides complete preview

---

## 9. Distribution

### 9.1 Package Managers

| Platform | Method | Command |
|----------|--------|---------|
| macOS | Homebrew | `brew install obay/tap/tfrn` |
| Windows | Scoop | `scoop bucket add obay https://github.com/obay/scoop-bucket && scoop install tfrn` |
| Linux | Homebrew | `brew install obay/tap/tfrn` |
| Any | Go | `go install github.com/obay/tfrn@latest` |

### 9.2 Release Artifacts

- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64

### 9.3 GoReleaser Configuration

GoReleaser will be configured to:
- Build `tfrn` binary for all target platforms
- Publish to GitHub Releases
- Update Homebrew tap
- Update Scoop bucket

---

## 10. Future Enhancements

### 10.1 Potential Features (v1.x)

- [ ] Custom naming patterns via configuration
- [ ] Watch mode for continuous organization
- [ ] VS Code extension

---

## 11. Appendix

### A. OTN Quick Reference

```
resource.azurerm_virtual_network.main.tf     # Resource
data.aws_ami.ubuntu-latest.tf                # Data source
provider.aws.tf                              # Provider
provider.aws.us-west-2.tf                    # Provider with alias
variable.environment.tf                      # Variable
output.vpc_id.tf                             # Output
module.networking.tf                         # Module
locals.tf                                    # Locals block
terraform.tf                                 # Terraform configuration
```

### B. Example Workflow

```bash
# Before
$ ls
main.tf  # Contains 10 resources, 3 data sources, 5 variables

# Run tfrn
$ tfrn organize --verbose

# After
$ ls
data.aws_ami.ubuntu.tf
data.aws_availability_zones.available.tf
data.aws_caller_identity.current.tf
locals.tf
module.vpc.tf
output.instance_id.tf
provider.aws.tf
resource.aws_instance.web.tf
resource.aws_security_group.web.tf
...
variable.environment.tf
variable.instance_type.tf
variable.region.tf
```

### C. JSON Schema (Output)

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "version": { "type": "string" },
    "command": { "type": "string", "enum": ["split", "rename", "organize", "validate"] },
    "directory": { "type": "string" },
    "dry_run": { "type": "boolean" },
    "results": {
      "type": "object",
      "properties": {
        "files_processed": { "type": "integer" },
        "files_split": { "type": "integer" },
        "files_renamed": { "type": "integer" },
        "files_skipped": { "type": "integer" },
        "files_already_compliant": { "type": "integer" },
        "errors": { "type": "array", "items": { "type": "string" } }
      }
    },
    "changes": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "action": { "type": "string", "enum": ["split", "rename", "skip"] },
          "source": { "type": "string" },
          "target": { "type": "string" },
          "targets": { "type": "array", "items": { "type": "string" } },
          "reason": { "type": "string" }
        }
      }
    },
    "violations": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "file": { "type": "string" },
          "current_name": { "type": "string" },
          "expected_name": { "type": "string" },
          "block_type": { "type": "string" }
        }
      }
    },
    "exit_code": { "type": "integer" }
  }
}
```

