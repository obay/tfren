# Story Map & Epic Breakdown

## tfrn - Terraform File Organizer

**Version:** 1.0
**Date:** December 2024
**Status:** Draft

---

## 1. Epic Overview

| Epic | Description | Priority |
|------|-------------|----------|
| E1 | Project Setup & CLI Framework | P0 - Critical |
| E2 | Core HCL Parsing | P0 - Critical |
| E3 | Naming Convention Engine | P0 - Critical |
| E4 | Split Operation | P0 - Critical |
| E5 | Rename Operation | P0 - Critical |
| E6 | Organize Operation | P0 - Critical |
| E7 | Validate Operation | P1 - High |
| E8 | Git Safety Checks | P1 - High |
| E9 | JSON Output | P1 - High |
| E10 | Configuration File Support | P2 - Medium |
| E11 | Release & Distribution | P1 - High |

---

## 2. User Stories by Epic

### Epic 1: Project Setup & CLI Framework

#### Story 1.1: Initialize Project Structure
**As a** developer
**I want** a well-organized project structure
**So that** the codebase is maintainable and follows Go best practices

**Acceptance Criteria:**
- [ ] Create directory structure as defined in TDD Section 2
- [ ] Initialize go.mod with module `github.com/obay/tfrn`
- [ ] Add Cobra and Viper dependencies
- [ ] Create placeholder files for all packages

**Tasks:**
- Create `cmd/tfrn/main.go`
- Create `internal/` package directories
- Run `go mod init` and `go mod tidy`

---

#### Story 1.2: Implement Root Command
**As a** user
**I want** to run `tfrn` with global flags
**So that** I can control the tool's behavior

**Acceptance Criteria:**
- [ ] `tfrn --version` prints version
- [ ] `tfrn --help` shows usage
- [ ] Global flags work: `-d`, `-r`, `-n`, `-v`, `-q`, `--json`
- [ ] Running `tfrn` without subcommand defaults to `organize`

**Tasks:**
- Implement `internal/cli/root.go` with Cobra
- Bind flags to Viper
- Add version command

---

#### Story 1.3: Implement Subcommand Stubs
**As a** user
**I want** subcommands for split, rename, organize, validate
**So that** I can perform specific operations

**Acceptance Criteria:**
- [ ] `tfrn split --help` shows split usage
- [ ] `tfrn rename --help` shows rename usage
- [ ] `tfrn organize --help` shows organize usage
- [ ] `tfrn validate --help` shows validate usage
- [ ] Each command accepts its specific flags

**Tasks:**
- Create `internal/cli/split.go`
- Create `internal/cli/rename.go`
- Create `internal/cli/organize.go`
- Create `internal/cli/validate.go`

---

### Epic 2: Core HCL Parsing

#### Story 2.1: Parse Single Terraform File
**As a** developer
**I want** to parse a .tf file into structured blocks
**So that** I can process each block independently

**Acceptance Criteria:**
- [ ] Parse valid HCL and return list of blocks
- [ ] Each block contains: type, labels, start/end lines
- [ ] Handle parse errors gracefully (return error, don't crash)
- [ ] Support all block types: resource, data, provider, variable, output, module, locals, terraform

**Tasks:**
- Implement `internal/hcl/parser.go`
- Implement `internal/hcl/block.go`
- Write unit tests for each block type

---

#### Story 2.2: Extract Comments with Blocks
**As a** user
**I want** comments before a block to stay with that block when splitting
**So that** my documentation is preserved

**Acceptance Criteria:**
- [ ] Contiguous comment lines before a block are captured
- [ ] Empty line breaks comment association
- [ ] File headers (comments before first block with empty line gap) are not captured
- [ ] Inline comments within blocks are preserved in raw content

**Tasks:**
- Implement `internal/hcl/comments.go`
- Implement `findPrecedingComments()` function
- Write unit tests with various comment patterns

---

#### Story 2.3: Extract Block Raw Content
**As a** developer
**I want** to extract the exact source text of a block including comments
**So that** split files preserve original formatting

**Acceptance Criteria:**
- [ ] Extract block content from source bytes using line ranges
- [ ] Include associated comments in extracted content
- [ ] Preserve original indentation and formatting

**Tasks:**
- Implement `extractRawContent()` function
- Handle edge cases (first block, last block, nested blocks)

---

### Epic 3: Naming Convention Engine

#### Story 3.1: Generate Filename from Block
**As a** developer
**I want** to generate OTN-compliant filenames from block metadata
**So that** files are named consistently

**Acceptance Criteria:**
- [ ] Resource: `resource.<provider>.<name>.tf`
- [ ] Data: `data.<provider>.<name>.tf`
- [ ] Provider: `provider.<name>.tf`
- [ ] Provider with alias: `provider.<name>.<alias>.tf`
- [ ] Variable: `variable.<name>.tf`
- [ ] Output: `output.<name>.tf`
- [ ] Module: `module.<name>.tf`
- [ ] Locals: `locals.tf`
- [ ] Terraform: `terraform.tf`
- [ ] Return empty string for unknown block types

**Tasks:**
- Implement `internal/naming/convention.go`
- Implement `GenerateFileName()` function
- Write unit tests for all block types

---

#### Story 3.2: Check Filename Compliance
**As a** developer
**I want** to check if a file's name matches its content
**So that** I can identify files needing rename

**Acceptance Criteria:**
- [ ] `IsCompliant(filename, block)` returns true if name matches convention
- [ ] Handle case sensitivity correctly

**Tasks:**
- Implement `IsCompliant()` function
- Write unit tests

---

### Epic 4: Split Operation

#### Story 4.1: Split Multi-Block File
**As a** user
**I want** to split a file with multiple blocks into separate files
**So that** each block has its own file

**Acceptance Criteria:**
- [ ] File with N blocks produces N new files
- [ ] Each new file is named according to OTN
- [ ] Original file is removed after successful split
- [ ] Comments are preserved with their blocks
- [ ] Handles file conflicts (target file already exists)

**Tasks:**
- Implement `internal/operations/split.go`
- Implement file writing logic
- Handle conflicts and errors

---

#### Story 4.2: Split with --keep-original Flag
**As a** user
**I want** to keep the original file after splitting
**So that** I can review changes before deleting

**Acceptance Criteria:**
- [ ] `--keep-original` flag prevents deletion of source file
- [ ] New files are still created

**Tasks:**
- Add flag handling in split command
- Pass option to split operation

---

#### Story 4.3: Split with --backup Flag
**As a** user
**I want** to create a backup before splitting
**So that** I can recover if something goes wrong

**Acceptance Criteria:**
- [ ] `--backup` creates `.bak` copy of original file
- [ ] Backup is created before any modifications

**Tasks:**
- Implement backup logic
- Add flag handling

---

#### Story 4.4: Split Dry Run
**As a** user
**I want** to preview what split would do
**So that** I can verify before making changes

**Acceptance Criteria:**
- [ ] `--dry-run` shows what files would be created
- [ ] No files are actually created or modified
- [ ] Output clearly indicates dry run mode

**Tasks:**
- Implement dry run path in split operation
- Return changes without executing

---

### Epic 5: Rename Operation

#### Story 5.1: Rename Single-Block File
**As a** user
**I want** to rename files to match OTN convention
**So that** my files are consistently named

**Acceptance Criteria:**
- [ ] Files with single block are renamed to OTN format
- [ ] Files already compliant are skipped with success message
- [ ] Files with multiple blocks are skipped with warning
- [ ] Handles file conflicts (target name already exists)

**Tasks:**
- Implement `internal/operations/rename.go`
- Handle edge cases and conflicts

---

#### Story 5.2: Rename with --backup Flag
**As a** user
**I want** to create a backup before renaming
**So that** I can recover the original filename

**Acceptance Criteria:**
- [ ] `--backup` creates `.bak` copy before rename
- [ ] Original content is preserved in backup

**Tasks:**
- Implement backup logic for rename
- Add flag handling

---

#### Story 5.3: Rename Dry Run
**As a** user
**I want** to preview renames before executing
**So that** I can verify the changes

**Acceptance Criteria:**
- [ ] `--dry-run` shows old → new filename mappings
- [ ] No files are actually renamed

**Tasks:**
- Implement dry run path in rename operation

---

### Epic 6: Organize Operation

#### Story 6.1: Combine Split and Rename
**As a** user
**I want** to run split and rename in one command
**So that** I can fully organize my Terraform files

**Acceptance Criteria:**
- [ ] First splits all multi-block files
- [ ] Then renames all files (including newly created ones)
- [ ] Reports summary of all changes
- [ ] `tfrn` without subcommand runs organize

**Tasks:**
- Implement `internal/operations/organize.go`
- Wire up as default command

---

#### Story 6.2: Organize Dry Run
**As a** user
**I want** to preview full organization
**So that** I can see the complete transformation

**Acceptance Criteria:**
- [ ] Shows all splits and renames that would occur
- [ ] No files are modified

**Tasks:**
- Combine dry run outputs from split and rename

---

### Epic 7: Validate Operation

#### Story 7.1: Check Naming Compliance
**As a** user
**I want** to check if files follow OTN convention
**So that** I can identify files needing attention

**Acceptance Criteria:**
- [ ] Lists all files that don't match convention
- [ ] Shows current name and expected name
- [ ] Reports count of compliant vs non-compliant files
- [ ] Exit code 0 if all compliant

**Tasks:**
- Implement `internal/operations/validate.go`
- Generate violation report

---

#### Story 7.2: Validate with --strict Flag
**As a** CI pipeline
**I want** validation to fail with non-zero exit code
**So that** I can enforce naming convention in CI

**Acceptance Criteria:**
- [ ] `--strict` returns exit code 2 if violations found
- [ ] Without `--strict`, always returns 0 (just reports)

**Tasks:**
- Add strict flag handling
- Return appropriate exit code

---

### Epic 8: Git Safety Checks

#### Story 8.1: Detect Git Repository
**As a** user
**I want** to be warned if I'm not in a Git repo
**So that** I know my changes can't be easily undone

**Acceptance Criteria:**
- [ ] Warning displayed if `.git` directory not found
- [ ] Warning can be suppressed with `--quiet`
- [ ] Warning does not block operation

**Tasks:**
- Implement `internal/git/status.go`
- Check for `.git` directory existence

---

#### Story 8.2: Detect Uncommitted Changes
**As a** user
**I want** to be warned if I have uncommitted changes
**So that** I can commit first and have a clean baseline

**Acceptance Criteria:**
- [ ] Warning displayed if `git status --porcelain` has output
- [ ] Warning can be suppressed with `--quiet`
- [ ] Warning does not block operation

**Tasks:**
- Run `git status --porcelain` and check output
- Display appropriate warning

---

### Epic 9: JSON Output

#### Story 9.1: JSON Result Structure
**As an** LLM or automation tool
**I want** structured JSON output
**So that** I can parse and act on results programmatically

**Acceptance Criteria:**
- [ ] `--json` flag outputs valid JSON to stdout
- [ ] JSON matches schema defined in PRD
- [ ] Includes: version, command, changes, violations, exit_code
- [ ] Human-readable output is suppressed when JSON is enabled

**Tasks:**
- Implement `internal/output/json.go`
- Define Result struct matching schema
- Integrate with all commands

---

#### Story 9.2: JSON Dry Run Output
**As an** LLM
**I want** dry run results in JSON
**So that** I can preview changes programmatically

**Acceptance Criteria:**
- [ ] `--json --dry-run` outputs planned changes
- [ ] `dry_run: true` is set in output

**Tasks:**
- Ensure dry run populates Result struct correctly

---

### Epic 10: Configuration File Support

#### Story 10.1: Load Configuration from File
**As a** user
**I want** to set defaults in a config file
**So that** I don't have to repeat flags every time

**Acceptance Criteria:**
- [ ] Reads `.tfrn.json` from current directory
- [ ] Reads `.tfrn.json` from home directory as fallback
- [ ] Command-line flags override config file values

**Tasks:**
- Implement `internal/config/config.go`
- Configure Viper to read JSON config
- Set up precedence chain

---

#### Story 10.2: Environment Variable Support
**As a** CI user
**I want** to configure via environment variables
**So that** I can set options in my CI environment

**Acceptance Criteria:**
- [ ] `TFRN_RECURSIVE=true` sets recursive mode
- [ ] `TFRN_DRY_RUN=true` sets dry run mode
- [ ] Environment variables override config file

**Tasks:**
- Configure Viper environment variable binding
- Use `TFRN_` prefix

---

### Epic 11: Release & Distribution

#### Story 11.1: GoReleaser Configuration
**As a** maintainer
**I want** automated releases
**So that** I can easily publish new versions

**Acceptance Criteria:**
- [ ] GoReleaser builds for linux/darwin/windows, amd64/arm64
- [ ] Creates GitHub release with binaries
- [ ] Updates Homebrew tap
- [ ] Updates Scoop bucket

**Tasks:**
- Update `.goreleaser.yaml`
- Configure Homebrew tap repository
- Configure Scoop bucket repository

---

#### Story 11.2: Version Injection
**As a** user
**I want** `tfrn --version` to show the correct version
**So that** I know which version I'm running

**Acceptance Criteria:**
- [ ] Version is injected at build time via ldflags
- [ ] Development builds show "dev"

**Tasks:**
- Update version.go
- Configure ldflags in GoReleaser

---

## 3. Story Dependencies

```
E1 (Project Setup)
 └── E2 (HCL Parsing)
      └── E3 (Naming)
           ├── E4 (Split)
           ├── E5 (Rename)
           │    └── E6 (Organize)
           └── E7 (Validate)

E8 (Git Safety) ─── Independent, integrates with E1
E9 (JSON Output) ── Independent, integrates with E4-E7
E10 (Config) ────── Independent, integrates with E1
E11 (Release) ───── After all others complete
```

---

## 4. Sprint Planning Suggestion

### Sprint 1: Foundation
- Story 1.1: Initialize Project Structure
- Story 1.2: Implement Root Command
- Story 1.3: Implement Subcommand Stubs
- Story 2.1: Parse Single Terraform File

### Sprint 2: Core Logic
- Story 2.2: Extract Comments with Blocks
- Story 2.3: Extract Block Raw Content
- Story 3.1: Generate Filename from Block
- Story 3.2: Check Filename Compliance

### Sprint 3: Operations
- Story 4.1: Split Multi-Block File
- Story 4.2: Split with --keep-original
- Story 4.4: Split Dry Run
- Story 5.1: Rename Single-Block File
- Story 5.3: Rename Dry Run

### Sprint 4: Polish
- Story 6.1: Combine Split and Rename
- Story 7.1: Check Naming Compliance
- Story 7.2: Validate with --strict
- Story 8.1: Detect Git Repository
- Story 8.2: Detect Uncommitted Changes

### Sprint 5: LLM & Release
- Story 9.1: JSON Result Structure
- Story 9.2: JSON Dry Run Output
- Story 10.1: Load Configuration from File
- Story 11.1: GoReleaser Configuration
- Story 11.2: Version Injection

---

## 5. Definition of Done

A story is complete when:
- [ ] Code is implemented and compiles
- [ ] Unit tests pass with >80% coverage for new code
- [ ] Integration tests pass (where applicable)
- [ ] Code is reviewed
- [ ] Documentation is updated (if user-facing)
- [ ] No known bugs or regressions

---

*Document Version: 1.0*
*Last Updated: December 2024*
