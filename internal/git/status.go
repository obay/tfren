package git

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Status int

const (
	NotARepo Status = iota
	CleanRepo
	DirtyRepo          // Has uncommitted changes (not .tf files in target dir)
	DirtyTerraformRepo // Has uncommitted .tf files in target dir - dangerous!
	GitNotInstalled
)

// CheckStatus checks git status for the target directory
// targetDir is the directory tfren will process
// recursive indicates if tfren will process subdirectories
func CheckStatus(targetDir string, recursive bool) Status {
	// Check if git is installed
	if _, err := exec.LookPath("git"); err != nil {
		return GitNotInstalled
	}

	// Check if we're in a git repo (works from any subfolder)
	cmd := exec.Command("git", "-C", targetDir, "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		return NotARepo
	}

	// Get uncommitted changes
	cmd = exec.Command("git", "-C", targetDir, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return NotARepo
	}

	changes := strings.TrimSpace(string(output))
	if changes == "" {
		return CleanRepo
	}

	// Get the absolute path of target directory for comparison
	absTargetDir, err := filepath.Abs(targetDir)
	if err != nil {
		return DirtyRepo
	}

	// Get git repo root
	cmd = exec.Command("git", "-C", targetDir, "rev-parse", "--show-toplevel")
	repoRootBytes, err := cmd.Output()
	if err != nil {
		return DirtyRepo
	}
	repoRoot := strings.TrimSpace(string(repoRootBytes))

	// Check if any uncommitted .tf files are in the target directory
	for _, line := range strings.Split(changes, "\n") {
		if len(line) < 4 {
			continue
		}
		// Git status format: "XY filename" or "XY original -> renamed"
		filePath := strings.TrimSpace(line[3:])
		// Handle renames: "old -> new"
		if idx := strings.Index(filePath, " -> "); idx != -1 {
			filePath = filePath[idx+4:]
		}

		if !strings.HasSuffix(filePath, ".tf") {
			continue
		}

		// Get absolute path of the changed file
		absFilePath := filepath.Join(repoRoot, filePath)

		// Check if file is in target directory
		if recursive {
			// If recursive, check if file is under target dir
			if strings.HasPrefix(absFilePath, absTargetDir+string(filepath.Separator)) || filepath.Dir(absFilePath) == absTargetDir {
				return DirtyTerraformRepo
			}
		} else {
			// If not recursive, check if file is directly in target dir
			if filepath.Dir(absFilePath) == absTargetDir {
				return DirtyTerraformRepo
			}
		}
	}

	return DirtyRepo
}

func (s Status) Warning() string {
	switch s {
	case GitNotInstalled:
		return "Note: Git is not installed. Consider using version control to track changes."
	case NotARepo:
		return "Warning: This directory is not a Git repository. Changes cannot be easily undone."
	case DirtyRepo:
		return "Warning: You have uncommitted changes."
	case DirtyTerraformRepo:
		return "Warning: You have uncommitted Terraform files in this directory. Changes made by tfren cannot be reverted with git."
	default:
		return ""
	}
}

func (s Status) RequiresConfirmation() bool {
	return s == DirtyTerraformRepo
}

// PromptContinue asks the user if they want to continue
func PromptContinue() bool {
	fmt.Print("Do you want to continue? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
