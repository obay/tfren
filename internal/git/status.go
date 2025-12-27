package git

import (
	"github.com/go-git/go-git/v5"
)

type Status int

const (
	NotARepo Status = iota
	DirtyRepo
	CleanRepo
)

func CheckStatus(dir string) Status {
	// Open repo, walking parent directories to find .git
	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return NotARepo
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return NotARepo
	}

	status, err := worktree.Status()
	if err != nil {
		return NotARepo
	}

	if status.IsClean() {
		return CleanRepo
	}
	return DirtyRepo
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
