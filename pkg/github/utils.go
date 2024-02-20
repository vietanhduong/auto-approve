package github

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func RepoRoot() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", fmt.Errorf("open git repository: %w", err)
	}
	tree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("get worktree: %w", err)
	}
	return tree.Filesystem.Root(), nil
}
