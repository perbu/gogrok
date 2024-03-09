package gitver

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"sort"
)

// GetLatestTag returns the latest tag in the repository
func GetLatestTag(repoPath string) (string, error) {
	gitRepo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", fmt.Errorf("git.PlainOpen: %w", err)

	}
	tagRefs, err := gitRepo.Tags()
	if err != nil {
		return "", fmt.Errorf("gitRepo.Tags: %w", err)
	}

	var tags []string
	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		tags = append(tags,
			t.Name().String(),
		)
		return nil
	})

	if len(tags) == 0 {
		return "", nil
	}

	// You might need to sort `tags` here based on your versioning scheme
	sort.Strings(tags)
	// Return the latest tag
	lTag := tags[len(tags)-1]
	// Remove the "refs/tags/" prefix
	return lTag[10:], nil
}
