package modver

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"golang.org/x/mod/semver"
)

// GetTags returns sorted a list of tags for a git repository
func fetchLocalTags(repoPath string) ([]string, error) {
	gitRepo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("git.PlainOpen: %w", err)
	}
	tagRefs, err := gitRepo.Tags()
	if err != nil {
		return nil, fmt.Errorf("gitRepo.Tags: %w", err)
	}

	var tags []string
	err = tagRefs.ForEach(func(t *plumbing.Reference) error {
		tags = append(tags,
			t.Name().String(),
		)
		return nil
	})

	if len(tags) == 0 {
		return nil, NoSuchRepoError{repoPath, "local, no tags"}
	}

	// remove the "refs/tags/" prefix from the tags
	for i, tag := range tags {
		tags[i] = tag[10:]
	}
	semver.Sort(tags)
	return tags, nil
}
