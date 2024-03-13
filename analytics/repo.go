package analytics

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

func New(path string) (*Repo, error) {
	r := &Repo{
		modules:  make(map[string]*Module),
		basePath: path,
	}
	return r, nil
}

func (r *Repo) Parse() error {
	start := time.Now()
	repoDirs, err := os.ReadDir(r.basePath)
	if err != nil {
		return fmt.Errorf("os.ReadDir: %w", err)
	}

	// First pass: parse all go.mod files:
	for _, repoDir := range repoDirs {
		if !repoDir.IsDir() {
			continue
		}
		if repoDir.Name() == "n2-metadata" {
			fmt.Println("Metadata found")
		}
		modulePath := path.Join(r.basePath, repoDir.Name())
		err := r.ParseMod(modulePath)
		if err != nil {
			return fmt.Errorf("r.ParseMod(%s): %w", modulePath, err)
		}
	}
	slog.Info("parsed go.mod files and git metadata", "duration", time.Since(start))
	start = time.Now()
	// Second pass: load and parse all source code:
	for _, name := range r.GetModuleNames() {
		mod, ok := r.GetModule(name)
		if !ok {
			return fmt.Errorf("r.GetModule(%s): not found", name)
		}
		if name == "github.com/celerway/n2-metadata" {
			fmt.Println("Metadata found")
		}
		// Load source code for local modules
		if mod.Type == DepTypeLocal {
			err := mod.LoadSource()
			if err != nil {
				return fmt.Errorf("mod.LoadSource: %w", err)
			}
			r.NoOfLines += mod.NoOfLines
			r.NoOfFiles += mod.NoOfFiles
		}
	}
	slog.Info("parsed source code", "duration", time.Since(start), "lines", r.NoOfLines, "files", r.NoOfFiles)
	// Populate reverse dependencies, both packages and modules.
	start = time.Now()
	r.reverseDeps()
	slog.Info("populated reverse dependencies", "duration", time.Since(start))
	return nil
}

func (r *Repo) ModuleFilter(t DepType, substring string) []Module {
	mods := make([]Module, 0)
	for _, v := range r.modules {
		if v.Type == t && (substring == "" || strings.Contains(v.Path, substring)) {
			mods = append(mods, *v)
		}
	}
	// Sort by module path
	sort.Slice(mods, func(i, j int) bool {
		return mods[i].Path < mods[j].Path
	})
	return mods
}
