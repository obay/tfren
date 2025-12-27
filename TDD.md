# Technical Design Document (TDD)

## tfrn - Terraform File Organizer

**Version:** 1.0
**Date:** December 2024
**Author:** Obay
**Status:** Draft

---

## 1. Overview

This document describes the technical architecture and implementation details for `tfrn`, a CLI tool for organizing Terraform configuration files according to Obay's Terraform Naming Convention (OTN).

### 1.1 Scope

This TDD covers:
- Project structure and package organization
- Core data structures and interfaces
- HCL parsing and comment handling
- CLI implementation with Cobra/Viper
- JSON output for LLM integration
- Git safety checks
- Error handling patterns
- Testing strategy

### 1.2 References

- [PRD.md](PRD.md) - Product Requirements Document
- [OTN Naming Convention](https://www.obay.cloud/post/terraform-naming-convention/)

---

## 2. Project Structure

```
tfrn/
├── cmd/
│   └── tfrn/
│       └── main.go                 # Entry point
├── internal/
│   ├── cli/
│   │   ├── root.go                 # Root command, global flags
│   │   ├── split.go                # Split subcommand
│   │   ├── rename.go               # Rename subcommand
│   │   ├── organize.go             # Organize subcommand
│   │   └── validate.go             # Validate subcommand
│   ├── config/
│   │   └── config.go               # Viper configuration handling
│   ├── git/
│   │   └── status.go               # Git repository status checks
│   ├── hcl/
│   │   ├── parser.go               # HCL parsing utilities
│   │   ├── block.go                # Block extraction and naming
│   │   └── comments.go             # Comment association logic
│   ├── naming/
│   │   └── convention.go           # OTN naming rules
│   ├── operations/
│   │   ├── split.go                # File splitting logic
│   │   ├── rename.go               # File renaming logic
│   │   ├── organize.go             # Combined split+rename
│   │   └── validate.go             # Validation logic
│   └── output/
│       ├── console.go              # Human-readable console output
│       └── json.go                 # JSON output for LLMs
├── go.mod
├── go.sum
└── version.go                      # Version information
```

---

## 3. Core Data Structures

### 3.1 Configuration

```go
// internal/config/config.go

type Config struct {
    Directory    string `json:"directory"`
    Recursive    bool   `json:"recursive"`
    DryRun       bool   `json:"dry_run"`
    Verbose      bool   `json:"verbose"`
    Quiet        bool   `json:"quiet"`
    JSONOutput   bool   `json:"json"`
}

func Load() (*Config, error)
func (c *Config) Validate() error
```

### 3.2 HCL Block Representation

```go
// internal/hcl/block.go

type Block struct {
    Type       string            // resource, data, variable, etc.
    Labels     []string          // e.g., ["aws_instance", "web"]
    Attributes map[string]string // For extracting alias, etc.
    StartLine  int               // Line number where block starts
    EndLine    int               // Line number where block ends
    Comments   []string          // Associated comments (lines before block)
    RawContent string            // Original source text including comments
}

type ParsedFile struct {
    Path       string
    Blocks     []Block
    ParseError error
}
```

### 3.3 Operation Results

```go
// internal/operations/result.go

type Action string

const (
    ActionSplit   Action = "split"
    ActionRename  Action = "rename"
    ActionSkip    Action = "skip"
)

type Change struct {
    Action  Action   `json:"action"`
    Source  string   `json:"source"`
    Target  string   `json:"target,omitempty"`
    Targets []string `json:"targets,omitempty"`
    Reason  string   `json:"reason,omitempty"`
}

type Violation struct {
    File         string `json:"file"`
    CurrentName  string `json:"current_name"`
    ExpectedName string `json:"expected_name"`
    BlockType    string `json:"block_type"`
}

type Result struct {
    Version              string      `json:"version"`
    Command              string      `json:"command"`
    Directory            string      `json:"directory"`
    DryRun               bool        `json:"dry_run"`
    FilesProcessed       int         `json:"files_processed"`
    FilesSplit           int         `json:"files_split"`
    FilesRenamed         int         `json:"files_renamed"`
    FilesSkipped         int         `json:"files_skipped"`
    FilesAlreadyCompliant int        `json:"files_already_compliant"`
    Changes              []Change    `json:"changes"`
    Violations           []Violation `json:"violations"`
    Errors               []string    `json:"errors"`
    ExitCode             int         `json:"exit_code"`
}
```

---

## 4. Component Design

### 4.1 CLI Layer (Cobra)

```go
// internal/cli/root.go

var rootCmd = &cobra.Command{
    Use:   "tfrn",
    Short: "Terraform file organizer",
    Long:  `tfrn organizes Terraform files according to OTN naming convention.`,
    PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
        // Load config, check git status, display warnings
        return preRun(cmd)
    },
}

func init() {
    rootCmd.PersistentFlags().StringP("directory", "d", ".", "Target directory")
    rootCmd.PersistentFlags().BoolP("recursive", "r", false, "Process subdirectories")
    rootCmd.PersistentFlags().BoolP("dry-run", "n", false, "Preview changes")
    rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")
    rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress warnings")
    rootCmd.PersistentFlags().Bool("json", false, "JSON output")

    // Bind to Viper
    viper.BindPFlags(rootCmd.PersistentFlags())
}
```

**Default Command Behavior:**
Running `tfrn` without a subcommand will execute `organize` (split + rename).

### 4.2 Git Safety Checks

```go
// internal/git/status.go

type GitStatus int

const (
    NotARepo GitStatus = iota
    DirtyRepo
    CleanRepo
)

func CheckStatus(dir string) GitStatus {
    // Check if .git exists
    // If exists, run `git status --porcelain`
    // Return appropriate status
}

func (s GitStatus) Warning() string {
    switch s {
    case NotARepo:
        return "Warning: This directory is not a Git repository. " +
               "Changes cannot be easily undone. Consider running 'git init' first."
    case DirtyRepo:
        return "Warning: You have uncommitted changes. " +
               "Consider committing or stashing before proceeding."
    default:
        return ""
    }
}
```

### 4.3 HCL Parsing with Comments

```go
// internal/hcl/parser.go

func ParseFile(path string) (*ParsedFile, error) {
    content, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    parser := hclparse.NewParser()
    file, diags := parser.ParseHCL(content, path)
    if diags.HasErrors() {
        return &ParsedFile{Path: path, ParseError: diags}, nil
    }

    blocks := extractBlocksWithComments(file, content)
    return &ParsedFile{Path: path, Blocks: blocks}, nil
}
```

```go
// internal/hcl/comments.go

func extractBlocksWithComments(file *hcl.File, source []byte) []Block {
    body := file.Body.(*hclsyntax.Body)
    lines := strings.Split(string(source), "\n")

    var blocks []Block
    for _, hclBlock := range body.Blocks {
        block := Block{
            Type:      hclBlock.Type,
            Labels:    hclBlock.Labels,
            StartLine: hclBlock.TypeRange.Start.Line,
            EndLine:   hclBlock.Body.(*hclsyntax.Body).EndRange.Line,
        }

        // Find associated comments (contiguous comment lines before block)
        block.Comments = findPrecedingComments(lines, block.StartLine)

        // Extract raw content including comments
        commentStart := block.StartLine - len(block.Comments) - 1
        if commentStart < 0 {
            commentStart = 0
        }
        block.RawContent = extractLines(lines, commentStart, block.EndLine)

        blocks = append(blocks, block)
    }

    return blocks
}

func findPrecedingComments(lines []string, blockStartLine int) []string {
    var comments []string

    // Walk backwards from the line before the block
    for i := blockStartLine - 2; i >= 0; i-- {
        line := strings.TrimSpace(lines[i])
        if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
            comments = append([]string{lines[i]}, comments...)
        } else if line == "" {
            // Empty line breaks comment association
            break
        } else {
            // Non-comment, non-empty line breaks association
            break
        }
    }

    return comments
}
```

### 4.4 Naming Convention

```go
// internal/naming/convention.go

func GenerateFileName(block hcl.Block) string {
    switch block.Type {
    case "resource", "data":
        if len(block.Labels) >= 2 {
            return fmt.Sprintf("%s.%s.%s.tf",
                block.Type, block.Labels[0], block.Labels[1])
        }
    case "provider":
        if len(block.Labels) >= 1 {
            name := block.Labels[0]
            if alias := block.Attributes["alias"]; alias != "" {
                return fmt.Sprintf("provider.%s.%s.tf", name, alias)
            }
            return fmt.Sprintf("provider.%s.tf", name)
        }
    case "variable", "module", "output":
        if len(block.Labels) >= 1 {
            return fmt.Sprintf("%s.%s.tf", block.Type, block.Labels[0])
        }
    case "locals":
        return "locals.tf"
    case "terraform":
        return "terraform.tf"
    }
    return ""
}

func IsCompliant(filename string, block hcl.Block) bool {
    expected := GenerateFileName(block)
    return filename == expected
}
```

### 4.5 Operations

```go
// internal/operations/split.go

type SplitOptions struct {
    KeepOriginal bool
    Backup       bool
    DryRun       bool
}

func Split(file string, opts SplitOptions) ([]Change, error) {
    parsed, err := hcl.ParseFile(file)
    if err != nil {
        return nil, err
    }

    if len(parsed.Blocks) <= 1 {
        return nil, nil // Nothing to split
    }

    var changes []Change
    for _, block := range parsed.Blocks {
        newName := naming.GenerateFileName(block)
        if newName == "" {
            continue
        }

        if !opts.DryRun {
            if err := writeBlockToFile(newName, block); err != nil {
                return changes, err
            }
        }

        changes = append(changes, Change{
            Action: ActionSplit,
            Source: file,
            Target: newName,
        })
    }

    if !opts.DryRun && !opts.KeepOriginal {
        os.Remove(file)
    }

    return changes, nil
}
```

```go
// internal/operations/rename.go

type RenameOptions struct {
    Backup bool
    DryRun bool
}

func Rename(file string, opts RenameOptions) (*Change, error) {
    parsed, err := hcl.ParseFile(file)
    if err != nil {
        return nil, err
    }

    if len(parsed.Blocks) != 1 {
        return &Change{
            Action: ActionSkip,
            Source: file,
            Reason: "file contains multiple blocks",
        }, nil
    }

    newName := naming.GenerateFileName(parsed.Blocks[0])
    if newName == "" || newName == file {
        return nil, nil // Already correct or can't determine name
    }

    if !opts.DryRun {
        if opts.Backup {
            copyFile(file, file+".bak")
        }
        if err := os.Rename(file, newName); err != nil {
            return nil, err
        }
    }

    return &Change{
        Action: ActionRename,
        Source: file,
        Target: newName,
    }, nil
}
```

### 4.6 Output Formatting

```go
// internal/output/console.go

type ConsoleOutput struct {
    verbose bool
    quiet   bool
}

func (c *ConsoleOutput) Success(msg string) {
    if !c.quiet {
        color.Green(msg)
    }
}

func (c *ConsoleOutput) Warning(msg string) {
    if !c.quiet {
        color.Yellow(msg)
    }
}

func (c *ConsoleOutput) Error(msg string) {
    color.Red(msg)
}

func (c *ConsoleOutput) Info(msg string) {
    if c.verbose && !c.quiet {
        fmt.Println(msg)
    }
}
```

```go
// internal/output/json.go

func PrintJSON(result *Result) error {
    encoder := json.NewEncoder(os.Stdout)
    encoder.SetIndent("", "  ")
    return encoder.Encode(result)
}
```

---

## 5. Comment Handling Algorithm

### 5.1 Rules

1. **Inline comments**: Preserved as part of the block's raw content
2. **Preceding comments**: Comments on contiguous lines immediately before a block are associated with that block
3. **Empty line**: An empty line breaks the association between comments and the following block
4. **File headers**: Comments at the top of the file (before any blocks) with an empty line separating them from the first block are discarded
5. **Floating comments**: Comments between blocks that are separated by empty lines from both blocks are discarded

### 5.2 Example

```hcl
# File header - DISCARDED
# Author: Obay

# This describes the VPC - KEPT (attached to next block)
resource "aws_vpc" "main" {
  cidr_block = "10.0.0.0/16"  # Inline - KEPT
}

# Floating comment - DISCARDED

# Subnet comment - KEPT
resource "aws_subnet" "public" {
  vpc_id = aws_vpc.main.id
}
```

### 5.3 Implementation

```go
func findPrecedingComments(lines []string, blockStartLine int) []string {
    var comments []string

    for i := blockStartLine - 2; i >= 0; i-- {
        line := strings.TrimSpace(lines[i])

        // Check if line is a comment
        if strings.HasPrefix(line, "#") ||
           strings.HasPrefix(line, "//") ||
           strings.HasPrefix(line, "/*") {
            comments = append([]string{lines[i]}, comments...)
        } else if line == "" {
            // Empty line breaks association
            break
        } else {
            // Non-comment content breaks association
            break
        }
    }

    return comments
}
```

---

## 6. Error Handling

### 6.1 Exit Codes

| Code | Constant | Meaning |
|------|----------|---------|
| 0 | `ExitSuccess` | All operations completed successfully |
| 1 | `ExitPartial` | Some files skipped (conflicts, parse errors) |
| 2 | `ExitValidationFailed` | Validation found violations (--strict) |
| 3 | `ExitFatal` | Fatal error (I/O, permissions) |

### 6.2 Error Types

```go
// internal/operations/errors.go

type FileConflictError struct {
    Source string
    Target string
}

func (e *FileConflictError) Error() string {
    return fmt.Sprintf("cannot rename %s to %s: file already exists",
        e.Source, e.Target)
}

type ParseError struct {
    File    string
    Message string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("failed to parse %s: %s", e.File, e.Message)
}
```

---

## 7. Testing Strategy

### 7.1 Unit Tests

| Package | Test Focus |
|---------|------------|
| `internal/naming` | Filename generation for all block types |
| `internal/hcl` | HCL parsing, comment extraction |
| `internal/git` | Git status detection |
| `internal/config` | Configuration loading and validation |

### 7.2 Integration Tests

```go
// internal/operations/split_test.go

func TestSplit_MultipleBlocks(t *testing.T) {
    // Create temp directory with test file
    // Run split operation
    // Verify correct files created with correct content
}

func TestSplit_PreservesComments(t *testing.T) {
    // Create file with block-level comments
    // Split and verify comments are attached to correct files
}

func TestSplit_DryRun(t *testing.T) {
    // Verify no files modified in dry-run mode
}
```

### 7.3 Test Fixtures

```
testdata/
├── single_block/
│   ├── input/
│   │   └── main.tf
│   └── expected/
│       └── resource.aws_instance.web.tf
├── multiple_blocks/
│   ├── input/
│   │   └── main.tf
│   └── expected/
│       ├── resource.aws_vpc.main.tf
│       ├── resource.aws_subnet.public.tf
│       └── data.aws_ami.ubuntu.tf
└── with_comments/
    ├── input/
    │   └── main.tf
    └── expected/
        ├── resource.aws_vpc.main.tf  # Should include comment
        └── resource.aws_subnet.public.tf
```

---

## 8. Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| `github.com/spf13/cobra` | v1.8+ | CLI framework |
| `github.com/spf13/viper` | v1.18+ | Configuration management |
| `github.com/hashicorp/hcl/v2` | v2.23+ | HCL parsing |
| `github.com/fatih/color` | v1.18+ | Terminal colors |
| `github.com/charmbracelet/log` | v0.4+ | Structured logging |

---

## 9. Build and Release

### 9.1 Build Command

```bash
go build -ldflags="-s -w -X main.Version=${VERSION}" -o tfrn ./cmd/tfrn
```

### 9.2 GoReleaser Configuration

Key points for `.goreleaser.yaml`:
- Build for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
- Homebrew tap: `obay/homebrew-tap`
- Scoop bucket: `obay/scoop-bucket`
- Binary name: `tfrn`

---

## 10. Migration from Existing Code

### 10.1 Current State

The existing codebase has:
- `main.go` - All logic in one file (~330 lines)
- `printhelpers.go` - Color output helpers
- `version.go` - Version variable
- Uses `flag` package for CLI

### 10.2 Migration Steps

1. Create new directory structure under `internal/`
2. Extract HCL parsing logic to `internal/hcl/`
3. Extract naming logic to `internal/naming/`
4. Extract file operations to `internal/operations/`
5. Implement Cobra CLI in `internal/cli/`
6. Add Viper configuration
7. Add Git status checks
8. Add JSON output support
9. Update entry point in `cmd/tfrn/main.go`
10. Add comprehensive tests

### 10.3 Backward Compatibility

The new CLI structure maintains compatibility:
- `tfrn` (no args) → runs `organize` (combines split + rename)
- `tfrn --version` → still works
- `tfrn split` → explicit split command
- `tfrn rename` → explicit rename command

---

## 11. Sign-Off

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Tech Lead | | | |
| Developer | | | |

---

*Document Version: 1.0*
*Last Updated: December 2024*
