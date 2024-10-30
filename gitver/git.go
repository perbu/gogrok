package gitver

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"golang.org/x/mod/semver"
	"regexp"
)

type NoSuchRepoError struct {
	repoPath string
}

func (e NoSuchRepoError) Error() string {
	return fmt.Sprintf("no such repo: '%s'", e.repoPath)
}

// GetLatestTag returns the latest tag in the (local) repository
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
	if err != nil {
		return "", fmt.Errorf("tagRefs.ForEach: %w", err)
	}
	if len(tags) == 0 {
		return "", nil
	}
	// You might need to sort `tags` here based on your versioning scheme
	semver.Sort(tags)
	// Return the latest tag
	lTag := tags[len(tags)-1]
	// Remove the "refs/tags/" prefix
	return lTag[10:], nil
}

func GetRemoteTags(url string) ([]string, error) {
	// Create a new in-memory storage object.
	storer := memory.NewStorage()
	rem, err := git.Init(storer, nil)
	if err != nil {
		return nil, fmt.Errorf("git.Init: %w", err)
	}

	// Add a remote with a name (e.g., "origin") and the URL.
	_, err = rem.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	})
	if err != nil {
		return nil, fmt.Errorf("rem.CreateRemote: %w", err)
	}

	// Fetch the tags from the remote.
	err = rem.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/tags/*:refs/tags/*"},
	})
	if err != nil {
		return nil, fmt.Errorf("rem.Fetch: %w", err)
	}

	// Get the list of tags from the remote.
	tags, err := rem.Tags()
	if err != nil {
		return nil, fmt.Errorf("rem.Tags: %w", err)
	}

	// make a list to store the tags:
	tagList := make([]string, 0, 10)
	// Iterate over the tags and print their names and hashes.
	err = tags.ForEach(func(t *plumbing.Reference) error {
		tagList = append(tagList, t.Name().Short())
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("tags.ForEach: %w", err)
	}
	// sort the tags
	semver.Sort(tagList)
	return tagList, nil
}

var pathRx = regexp.MustCompile(`^/([^/]+)/(.*)$`)

func path2url(path string) (string, error) {
	// check if the path is valid and extract the host
	m := pathRx.FindStringSubmatch(path)
	if m == nil {
		return "", fmt.Errorf("invalid path: '%s'", path)
	}
	if len(m) != 3 {
		return "", fmt.Errorf("invalid path: '%s'", path)
	}
	host := m[1]
	var urlStr string
	switch host {
	case "github.com":
		urlStr = fmt.Sprintf("https://%s", path)
	case "gitlab.com":
		urlStr = fmt.Sprintf("https://%s", path)
	default:
		return "", fmt.Errorf("unsupported host: '%s'", host)
	}
	return urlStr, nil
}
