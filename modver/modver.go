package modver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"golang.org/x/mod/semver"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

const (
	versionListQuery = "https://proxy.golang.org/%s/@v/list"
	dbFileName       = ".cache.bolt.db"
)

var httpClient = &http.Client{Timeout: 5 * time.Second}

type NoSuchRepoError struct {
	repoPath   string
	additional string
}

func (e NoSuchRepoError) Error() string {
	return fmt.Sprintf("no such repo or no versions(%s) : '%s'", e.additional, e.repoPath)
}

type ModTracker struct {
	cache  *bolt.DB
	logger *slog.Logger
	// cache is a map of module paths to their versions
}

func New() *ModTracker {
	db, err := bolt.Open(dbFileName, 0600, nil)
	if err != nil {
		panic(err)
	}
	return &ModTracker{
		cache: db,
	}
}

func (m *ModTracker) Close() error {
	if m.cache == nil {
		return errors.New("cache is already closed")
	}
	err := m.cache.Close()
	if err != nil {
		m.logger.Error("failed to close bolt db", "error", err)
	}
	m.cache = nil
	return err
}

type Module struct {
	Versions []string  `json:"versions"`
	Ts       time.Time `json:"ts"`
}

func (m *ModTracker) GetTags(ctx context.Context, module string) ([]string, error) {
	// 1. db is already opened in New()

	// 2. Lowercase the module path
	module = strings.ToLower(module)

	// 3. Attempt to fetch from cache
	var tags []string
	err := m.cache.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("modules"))
		if b == nil {
			return nil // Bucket doesn't exist, treat as cache miss
		}

		cachedData := b.Get([]byte(module))
		if cachedData == nil {
			return nil // Cache miss
		}

		var modules Module
		err := json.Unmarshal(cachedData, &modules)
		if err != nil {
			// return fmt.Errorf("failed to unmarshal cached data: %w", err)
			return nil // Treat as cache miss.
		}

		// Check if cache entry is expired
		if time.Since(modules.Ts) > 24*time.Hour {
			return nil // Cache expired
		}

		tags = modules.Versions
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read from cache: %w", err)
	}

	// 4. Cache hit: return cached tags
	if tags != nil {
		return tags, nil
	}
	// 5. Cache miss: fetch tags (your existing logic)
	var fetchErr error
	if !strings.HasPrefix(module, "code/") {
		tags, fetchErr = fetchRemoteTags(ctx, module)
	} else {
		tags, fetchErr = fetchLocalTags(module)
	}

	if fetchErr != nil {
		// Cache the NoSuchRepoError
		var noSuchRepoError NoSuchRepoError
		if errors.As(fetchErr, &noSuchRepoError) {
			err = m.cache.Update(func(tx *bolt.Tx) error {
				b, err := tx.CreateBucketIfNotExists([]byte("modules"))
				if err != nil {
					return err
				}
				return b.Put([]byte(module), nil) // Store nil to indicate error
			})
			if err != nil {
				return nil, fmt.Errorf("failed to cache error: %w", err)
			}
			return nil, fetchErr
		}
		return nil, fmt.Errorf("failed to fetch tags: %w", fetchErr)
	}

	// 6. Store fetched tags in cache
	err = m.cache.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("modules"))
		if err != nil {
			return err
		}

		modules := Module{
			Versions: tags,
			Ts:       time.Now(),
		}
		cachedData, err := json.Marshal(modules)
		if err != nil {
			return fmt.Errorf("failed to marshal data for cache: %w", err)
		}

		return b.Put([]byte(module), cachedData)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to write to cache: %w", err)
	}

	return tags, nil
}

func fetchRemoteTags(ctx context.Context, module string) ([]string, error) {
	// cache miss:
	u := fmt.Sprintf(versionListQuery, module)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http.Client.Do(req): %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, NoSuchRepoError{module, "remote not found"}
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	// there is one version on each line of the body response
	// split the body by new line using the default split function:
	tagList := strings.Split(string(body), "\n")
	// remove the empty strings:
	tagList = tagList[:len(tagList)-1] // remove the last empty string
	// sort the tags
	semver.Sort(tagList)
	return tagList, nil
}

func (m *ModTracker) GetLatestVersion(ctx context.Context, module string) (string, error) {
	tags, err := m.GetTags(ctx, module)
	if err != nil {
		return "", fmt.Errorf("GetTags(%s): %w", module, err)
	}
	if len(tags) == 0 {
		return "", NoSuchRepoError{module, "cache, no tags"}
	}
	return tags[len(tags)-1], nil // the last tag is the latest, list is sorted
}
