package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Status int

const (
	NotARepo Status = iota
	DirtyRepo
	CleanRepo
)

func CheckStatus(dir string) Status {
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return NotARepo
	}

	cmd := exec.Command("git", "-C", dir, "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return NotARepo
	}

	if strings.TrimSpace(string(output)) != "" {
		return DirtyRepo
	}

	return CleanRepo
}

func (s Status) Warning() string {
	switch s {
	case NotARepo:
		return "Warning: This directory is not a Git repository. Changes cannot be easily undone. Consider running 'git init' first."
	case DirtyRepo:
		return "Warning: You have uncommitted changes. Consider committing or stashing before proceeding."
	default:
		return ""
	}
}
